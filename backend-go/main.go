// OMNIXIUS API — Go backend. Stack: Go only (per project policy).
package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"omnixius-api/db"
	"omnixius-api/pqc"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

var cfg Config

// per-IP login rate limit: 5 attempts per 15 min
var loginLimitMu sync.Mutex
var loginLimiters = make(map[string]*rate.Limiter)

func main() {
	cfg = LoadConfig()
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("db", "omnixius.db")
	}
	if err := db.Open(dbPath); err != nil {
		panic(err)
	}
	db.InitUploadDirs(cfg.UploadDir)

	r := gin.New()
	r.Use(requestLogger())
	r.Use(corsMiddleware())
	r.Use(rateLimitMiddleware())
	r.Static("/uploads", cfg.UploadDir)

	api := r.Group("/api")
	api.POST("/auth/register", handleRegister)
	api.POST("/auth/login", handleLogin)
	api.GET("/auth/confirm-email", handleConfirmEmail)
	api.POST("/auth/forgot-password", handleForgotPassword)
	api.POST("/auth/reset-password", handleResetPassword)

	auth := api.Group("")
	auth.Use(authRequired())
	auth.GET("/users/me", handleUserMe)
	auth.PATCH("/users/me", handleUserUpdate)
	auth.GET("/users/me/orders", handleUserOrders)
	auth.POST("/users/me/avatar", handleUserAvatar)

	api.GET("/products", handleProductsList)
	api.GET("/products/categories", handleProductsCategories)
	api.GET("/products/:id", handleProductGet)
	auth.POST("/products", handleProductCreate)
	auth.PATCH("/products/:id", handleProductUpdate)
	auth.DELETE("/products/:id", handleProductDelete)

	auth.GET("/orders/my", handleOrdersMy)
	auth.POST("/orders", handleOrderCreate)
	auth.PATCH("/orders/:id", handleOrderUpdate)

	auth.GET("/conversations", handleConversationsList)
	auth.POST("/conversations", handleConversationCreate)
	auth.GET("/messages/conversation/:id", handleMessagesList)
	auth.POST("/messages/conversation/:id", handleMessageSend)
	auth.POST("/messages/:id/read", handleMessageRead)

	port := ":" + cfg.Port
	if err := r.Run(port); err != nil {
		panic(err)
	}
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		c.Next()
		status := c.Writer.Status()
		latency := time.Since(start)
		log.Printf("[%s] %d %s %s %s", method, status, path, clientIP, latency)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rate.NewLimiter(rate.Every(time.Minute/200), 200)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tok := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tok == "" {
			c.JSON(401, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}
		uid, _, err := pqc.VerifyToken(cfg.PQCPublicKey, tok)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		var email, role string
		var name, avatar sql.NullString
		var verified int
		err = db.DB.QueryRow("SELECT id, email, role, name, avatar_path, email_verified FROM users WHERE id = ?", uid).Scan(
			&uid, &email, &role, &name, &avatar, &verified)
		if err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		c.Set("userID", uid)
		c.Set("userName", name)
		c.Set("userAvatar", avatar)
		c.Next()
	}
}

func getUserID(c *gin.Context) int64 {
	v, _ := c.Get("userID")
	if id, ok := v.(int64); ok {
		return id
	}
	return 0
}

func getLoginLimiter(ip string) *rate.Limiter {
	loginLimitMu.Lock()
	defer loginLimitMu.Unlock()
	if l, ok := loginLimiters[ip]; ok {
		return l
	}
	l := rate.NewLimiter(rate.Every(15*time.Minute/5), 5) // 5 per 15 min
	loginLimiters[ip] = l
	return l
}

func handleRegister(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Email == "" || len(body.Password) < 8 {
		c.JSON(400, gin.H{"error": "Invalid email or password (min 8 chars)"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if len(body.Email) > 255 {
		c.JSON(400, gin.H{"error": "Email too long"})
		return
	}
	if len(body.Name) > 200 {
		body.Name = body.Name[:200]
	}
	var id int64
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = ?", body.Email).Scan(&id)
	if err == nil {
		c.JSON(409, gin.H{"error": "Email already registered"})
		return
	}
	hash := hashPasswordArgon2(body.Password)
	tok := make([]byte, 32)
	rand.Read(tok)
	verifyToken := hex.EncodeToString(tok)
	res, err := db.DB.Exec("INSERT INTO users (email, password_hash, name, role, email_verify_token) VALUES (?, ?, ?, 'user', ?)",
		body.Email, hash, nullStr(body.Name), verifyToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Registration failed"})
		return
	}
	id, _ = res.LastInsertId()
	token, _ := pqc.SignToken(cfg.PQCPrivateKey, id, time.Now().Add(7*24*time.Hour))
	c.JSON(201, gin.H{"user": gin.H{"id": id, "email": body.Email, "role": "user", "name": body.Name}, "token": token})
}

func handleLogin(c *gin.Context) {
	ip := c.ClientIP()
	limiter := getLoginLimiter(ip)
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many login attempts. Try again later."})
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Email and password required"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	var id int64
	var hash string
	var role, name string
	var avatar sql.NullString
	var verified int
	err := db.DB.QueryRow("SELECT id, password_hash, role, name, avatar_path, email_verified FROM users WHERE email = ?", body.Email).
		Scan(&id, &hash, &role, &name, &avatar, &verified)
	if err != nil || !checkPassword(hash, body.Password) {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}
	user := gin.H{"id": id, "email": body.Email, "role": role, "name": name, "avatar_path": avatar.String, "email_verified": verified == 1}
	token, _ := pqc.SignToken(cfg.PQCPrivateKey, id, time.Now().Add(7*24*time.Hour))
	c.JSON(200, gin.H{"user": user, "token": token})
}

func handleConfirmEmail(c *gin.Context) {
	tok := c.Query("token")
	if tok == "" {
		c.JSON(400, gin.H{"error": "Token required"})
		return
	}
	res, err := db.DB.Exec("UPDATE users SET email_verified = 1, email_verify_token = NULL WHERE email_verify_token = ?", tok)
	if err != nil || mustRows(res) == 0 {
		c.JSON(400, gin.H{"error": "Invalid token"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func handleForgotPassword(c *gin.Context) {
	var body struct{ Email string `json:"email"` }
	c.ShouldBindJSON(&body)
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if body.Email == "" {
		c.JSON(400, gin.H{"error": "Email required"})
		return
	}
	var id int64
	if db.DB.QueryRow("SELECT id FROM users WHERE email = ?", body.Email).Scan(&id) != nil {
		c.JSON(200, gin.H{"ok": true})
		return
	}
	tok := make([]byte, 32)
	rand.Read(tok)
	resetToken := hex.EncodeToString(tok)
	exp := time.Now().Add(time.Hour).Unix()
	db.DB.Exec("UPDATE users SET reset_token = ?, reset_token_expires = ? WHERE id = ?", resetToken, exp, id)
	c.JSON(200, gin.H{"ok": true})
}

func handleResetPassword(c *gin.Context) {
	var body struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Token == "" || len(body.Password) < 8 {
		c.JSON(400, gin.H{"error": "Token and password (min 8 chars) required"})
		return
	}
	var id int64
	err := db.DB.QueryRow("SELECT id FROM users WHERE reset_token = ? AND reset_token_expires > ?", body.Token, time.Now().Unix()).Scan(&id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid or expired token"})
		return
	}
	db.DB.Exec("UPDATE users SET password_hash = ?, reset_token = NULL, reset_token_expires = NULL WHERE id = ?", hashPasswordArgon2(body.Password), id)
	c.JSON(200, gin.H{"ok": true})
}

const argon2SaltLen = 16
const argon2HashLen = 32

func hashPasswordArgon2(password string) string {
	salt := make([]byte, argon2SaltLen)
	rand.Read(salt)
	key := argon2.IDKey([]byte(password), salt, cfg.Argon2Time, cfg.Argon2Memory, cfg.Argon2Threads, argon2HashLen)
	return "argon2id:" + base64.RawStdEncoding.EncodeToString(append(salt, key...))
}

func checkPassword(stored, password string) bool {
	if strings.HasPrefix(stored, "$2") {
		return checkPasswordBcrypt(stored, password)
	}
	if !strings.HasPrefix(stored, "argon2id:") {
		return false
	}
	b, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(stored, "argon2id:"))
	if err != nil || len(b) != argon2SaltLen+argon2HashLen {
		return false
	}
	salt := b[:argon2SaltLen]
	want := b[argon2SaltLen:]
	got := argon2.IDKey([]byte(password), salt, cfg.Argon2Time, cfg.Argon2Memory, cfg.Argon2Threads, argon2HashLen)
	return subtle.ConstantTimeCompare(want, got) == 1
}

func checkPasswordBcrypt(stored, password string) bool {
	// Legacy bcrypt (keep dependency for old hashes only)
	return bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) == nil
}

func handleUserMe(c *gin.Context) {
	id := getUserID(c)
	var email, role, name string
	var avatar sql.NullString
	var verified int
	if db.DB.QueryRow("SELECT email, role, name, avatar_path, email_verified FROM users WHERE id = ?", id).Scan(&email, &role, &name, &avatar, &verified) != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{"id": id, "email": email, "role": role, "name": name, "avatar_path": avatar.String, "email_verified": verified == 1})
}

func handleUserUpdate(c *gin.Context) {
	var body struct{ Name string `json:"name"` }
	c.ShouldBindJSON(&body)
	if len(body.Name) > 200 {
		body.Name = body.Name[:200]
	}
	db.DB.Exec("UPDATE users SET name = ?, updated_at = unixepoch() WHERE id = ?", nullStr(body.Name), getUserID(c))
	handleUserMe(c)
}

func handleUserOrders(c *gin.Context) {
	id := getUserID(c)
	// asBuyer
	rows, _ := db.DB.Query(`SELECT o.id, o.status, o.created_at, p.title, p.price, p.image_path, u.name FROM orders o JOIN products p ON p.id = o.product_id JOIN users u ON u.id = o.seller_id WHERE o.buyer_id = ? ORDER BY o.created_at DESC`, id)
	asBuyer := rowsToOrderList(rows)
	rows, _ = db.DB.Query(`SELECT o.id, o.status, o.created_at, p.title, p.price, p.image_path, u.name FROM orders o JOIN products p ON p.id = o.product_id JOIN users u ON u.id = o.buyer_id WHERE o.seller_id = ? ORDER BY o.created_at DESC`, id)
	asSeller := rowsToOrderList(rows)
	c.JSON(200, gin.H{"asBuyer": asBuyer, "asSeller": asSeller})
}

func rowsToOrderList(rows *sql.Rows) []gin.H {
	var out []gin.H
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, created int64
		var status, title string
		var price float64
		var imagePath, name sql.NullString
		rows.Scan(&id, &status, &created, &title, &price, &imagePath, &name)
		out = append(out, gin.H{"id": id, "status": status, "created_at": created, "title": title, "price": price, "image_path": imagePath.String, "seller_name": name.String})
	}
	return out
}

func handleUserAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, gin.H{"error": "File required"})
		return
	}
	ext := ".jpg"
	if strings.Contains(file.Header.Get("Content-Type"), "png") {
		ext = ".png"
	} else if strings.Contains(file.Header.Get("Content-Type"), "webp") {
		ext = ".webp"
	}
	rel := filepath.Join("avatars", strconv.FormatInt(time.Now().Unix(), 10)+ext)
	dst := filepath.Join(cfg.UploadDir, rel)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"error": "Upload failed"})
		return
	}
	db.DB.Exec("UPDATE users SET avatar_path = ?, updated_at = unixepoch() WHERE id = ?", rel, getUserID(c))
	c.JSON(200, gin.H{"avatar_path": rel})
}

func handleProductsList(c *gin.Context) {
	q := c.Query("q")
	cat := c.Query("category")
	loc := c.Query("location")
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	offset := 0
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	qry := `SELECT p.id, p.title, p.price, p.category, p.location, p.image_path, p.created_at, u.id, u.name FROM products p JOIN users u ON u.id = p.user_id WHERE 1=1`
	args := []interface{}{}
	if q != "" {
		qry += ` AND (p.title LIKE ? OR p.description LIKE ?)`
		args = append(args, "%"+q+"%", "%"+q+"%")
	}
	if cat != "" {
		qry += ` AND p.category = ?`
		args = append(args, cat)
	}
	if loc != "" {
		qry += ` AND p.location LIKE ?`
		args = append(args, "%"+loc+"%")
	}
	if v := c.Query("minPrice"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			qry += ` AND p.price >= ?`
			args = append(args, f)
		}
	}
	if v := c.Query("maxPrice"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			qry += ` AND p.price <= ?`
			args = append(args, f)
		}
	}
	qry += ` ORDER BY p.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, _ := db.DB.Query(qry, args...)
	var list []gin.H
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var id, created int64
			var title, category, location string
			var price float64
			var imagePath, sellerName sql.NullString
			var sellerID int64
			rows.Scan(&id, &title, &price, &category, &location, &imagePath, &created, &sellerID, &sellerName)
			list = append(list, gin.H{"id": id, "title": title, "price": price, "category": category, "location": location, "image_path": imagePath.String, "created_at": created, "seller_id": sellerID, "seller_name": sellerName.String})
		}
	}
	c.JSON(200, list)
}

func handleProductsCategories(c *gin.Context) {
	rows, _ := db.DB.Query("SELECT DISTINCT category FROM products ORDER BY category")
	var list []string
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var s string
			rows.Scan(&s)
			list = append(list, s)
		}
	}
	c.JSON(200, list)
}

func handleProductGet(c *gin.Context) {
	id := c.Param("id")
	var p productRow
	err := db.DB.QueryRow(`SELECT p.id, p.user_id, p.title, p.description, p.price, p.category, p.location, p.image_path, p.created_at, u.name, u.email FROM products p JOIN users u ON u.id = p.user_id WHERE p.id = ?`, id).
		Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &p.Price, &p.Category, &p.Location, &p.ImagePath, &p.CreatedAt, &p.SellerName, &p.SellerEmail)
	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(200, p.toH())
}

type productRow struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Price       float64
	Category    string
	Location    string
	ImagePath   sql.NullString
	CreatedAt   int64
	SellerName  sql.NullString
	SellerEmail string
}

func (p productRow) toH() gin.H {
	return gin.H{"id": p.ID, "user_id": p.UserID, "title": p.Title, "description": p.Description, "price": p.Price, "category": p.Category, "location": p.Location, "image_path": p.ImagePath.String, "created_at": p.CreatedAt, "seller_id": p.UserID, "seller_name": p.SellerName.String, "seller_email": p.SellerEmail}
}

func handleProductCreate(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	desc := strings.TrimSpace(c.PostForm("description"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	category := strings.TrimSpace(c.PostForm("category"))
	location := strings.TrimSpace(c.PostForm("location"))
	if title == "" || len(title) < 2 || category == "" || price < 0 {
		c.JSON(400, gin.H{"error": "Title (min 2), category and non-negative price required"})
		return
	}
	if len(title) > 200 {
		title = title[:200]
	}
	if len(desc) > 5000 {
		desc = desc[:5000]
	}
	if len(category) > 100 {
		category = category[:100]
	}
	if len(location) > 200 {
		location = location[:200]
	}
	imagePath := ""
	if file, err := c.FormFile("image"); err == nil {
		ext := ".jpg"
		if strings.Contains(file.Header.Get("Content-Type"), "png") {
			ext = ".png"
		} else if strings.Contains(file.Header.Get("Content-Type"), "webp") {
			ext = ".webp"
		}
		rel := filepath.Join("products", strconv.FormatInt(time.Now().Unix(), 10)+ext)
		dst := filepath.Join(cfg.UploadDir, rel)
		if c.SaveUploadedFile(file, dst) == nil {
			imagePath = rel
		}
	}
	res, err := db.DB.Exec("INSERT INTO products (user_id, title, description, price, category, location, image_path) VALUES (?, ?, ?, ?, ?, ?, ?)",
		getUserID(c), title, desc, price, category, nullStr(location), nullStr(imagePath))
	if err != nil {
		c.JSON(500, gin.H{"error": "Create failed"})
		return
	}
	newID, _ := res.LastInsertId()
	row := db.DB.QueryRow("SELECT id, user_id, title, description, price, category, location, image_path, created_at FROM products WHERE id = ?", newID)
	var pr productRow
	row.Scan(&pr.ID, &pr.UserID, &pr.Title, &pr.Description, &pr.Price, &pr.Category, &pr.Location, &pr.ImagePath, &pr.CreatedAt)
	c.JSON(201, pr.toH())
}

func handleProductUpdate(c *gin.Context) {
	id := c.Param("id")
	var ownerID int64
	if db.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", id).Scan(&ownerID) != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	if ownerID != getUserID(c) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}
	title := strings.TrimSpace(c.PostForm("title"))
	desc := strings.TrimSpace(c.PostForm("description"))
	category := strings.TrimSpace(c.PostForm("category"))
	location := strings.TrimSpace(c.PostForm("location"))
	if title != "" {
		db.DB.Exec("UPDATE products SET title = ? WHERE id = ?", title, id)
	}
	if desc != "" || c.PostForm("description") != "" {
		db.DB.Exec("UPDATE products SET description = ? WHERE id = ?", desc, id)
	}
	if v := c.PostForm("price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			db.DB.Exec("UPDATE products SET price = ? WHERE id = ?", f, id)
		}
	}
	if category != "" {
		db.DB.Exec("UPDATE products SET category = ? WHERE id = ?", category, id)
	}
	db.DB.Exec("UPDATE products SET location = ?, updated_at = unixepoch() WHERE id = ?", nullStr(location), id)
	if file, err := c.FormFile("image"); err == nil {
		rel := filepath.Join("products", strconv.FormatInt(time.Now().Unix(), 10)+".jpg")
		dst := filepath.Join(cfg.UploadDir, rel)
		if c.SaveUploadedFile(file, dst) == nil {
			db.DB.Exec("UPDATE products SET image_path = ? WHERE id = ?", rel, id)
		}
	}
	handleProductGet(c)
}

func handleProductDelete(c *gin.Context) {
	id := c.Param("id")
	var ownerID int64
	if db.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", id).Scan(&ownerID) != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	if ownerID != getUserID(c) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}
	db.DB.Exec("DELETE FROM products WHERE id = ?", id)
	c.Status(204)
}

func handleOrdersMy(c *gin.Context) {
	id := getUserID(c)
	rows, _ := db.DB.Query(`SELECT o.id, o.product_id, o.buyer_id, o.seller_id, o.status, o.created_at, p.title, p.price, p.image_path FROM orders o JOIN products p ON p.id = o.product_id WHERE o.buyer_id = ? OR o.seller_id = ? ORDER BY o.created_at DESC`, id, id)
	var list []gin.H
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var oid, pid, buyer, seller, created int64
			var status, title string
			var price float64
			var img sql.NullString
			rows.Scan(&oid, &pid, &buyer, &seller, &status, &created, &title, &price, &img)
			list = append(list, gin.H{"id": oid, "product_id": pid, "buyer_id": buyer, "seller_id": seller, "status": status, "created_at": created, "title": title, "price": price, "image_path": img.String})
		}
	}
	c.JSON(200, list)
}

func handleOrderCreate(c *gin.Context) {
	var body struct {
		ProductID int64 `json:"product_id"`
	}
	if c.ShouldBindJSON(&body) != nil || body.ProductID == 0 {
		c.JSON(400, gin.H{"error": "product_id required"})
		return
	}
	var sellerID int64
	if db.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", body.ProductID).Scan(&sellerID) != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	uid := getUserID(c)
	if sellerID == uid {
		c.JSON(400, gin.H{"error": "Cannot order own product"})
		return
	}
	res, err := db.DB.Exec("INSERT INTO orders (product_id, buyer_id, seller_id, status) VALUES (?, ?, ?, 'pending')", body.ProductID, uid, sellerID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	oid, _ := res.LastInsertId()
	c.JSON(201, gin.H{"id": oid, "product_id": body.ProductID, "buyer_id": uid, "seller_id": sellerID, "status": "pending"})
}

func handleOrderUpdate(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status"`
	}
	c.ShouldBindJSON(&body)
	if body.Status != "pending" && body.Status != "confirmed" && body.Status != "completed" && body.Status != "cancelled" {
		c.JSON(400, gin.H{"error": "Invalid status"})
		return
	}
	var buyer, seller int64
	if db.DB.QueryRow("SELECT buyer_id, seller_id FROM orders WHERE id = ?", id).Scan(&buyer, &seller) != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}
	uid := getUserID(c)
	if buyer != uid && seller != uid {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}
	db.DB.Exec("UPDATE orders SET status = ?, updated_at = unixepoch() WHERE id = ?", body.Status, id)
	c.JSON(200, gin.H{"id": id, "status": body.Status})
}

func handleConversationsList(c *gin.Context) {
	uid := getUserID(c)
	rows, _ := db.DB.Query(`SELECT c.id, c.product_id, c.updated_at FROM conversations c JOIN conversation_participants cp ON cp.conversation_id = c.id WHERE cp.user_id = ? ORDER BY c.updated_at DESC`, uid)
	var list []gin.H
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var cid, pid, updated int64
			var pidNull sql.NullInt64
			rows.Scan(&cid, &pidNull, &updated)
			if pidNull.Valid {
				pid = pidNull.Int64
			}
			var otherID int64
			var otherName, otherEmail sql.NullString
			db.DB.QueryRow("SELECT u.id, u.name, u.email FROM users u JOIN conversation_participants cp ON cp.user_id = u.id WHERE cp.conversation_id = ? AND u.id != ?", cid, uid).Scan(&otherID, &otherName, &otherEmail)
			var lastMsg sql.NullString
			db.DB.QueryRow("SELECT body FROM messages WHERE conversation_id = ? ORDER BY created_at DESC LIMIT 1", cid).Scan(&lastMsg)
			list = append(list, gin.H{"id": cid, "product_id": pid, "updated_at": updated, "last_message": lastMsg.String, "other": gin.H{"id": otherID, "name": otherName.String, "email": otherEmail.String}})
		}
	}
	c.JSON(200, list)
}

func handleConversationCreate(c *gin.Context) {
	var body struct {
		UserID    int64 `json:"user_id"`
		ProductID int64 `json:"product_id"`
	}
	c.ShouldBindJSON(&body)
	if body.UserID == 0 || body.UserID == getUserID(c) {
		c.JSON(400, gin.H{"error": "Valid user_id required"})
		return
	}
	var otherExists int64
	if db.DB.QueryRow("SELECT id FROM users WHERE id = ?", body.UserID).Scan(&otherExists) != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	uid := getUserID(c)
	var convID int64
	err := db.DB.QueryRow(`SELECT c.id FROM conversations c JOIN conversation_participants cp1 ON cp1.conversation_id = c.id AND cp1.user_id = ? JOIN conversation_participants cp2 ON cp2.conversation_id = c.id AND cp2.user_id = ? WHERE (c.product_id IS NULL AND ? = 0) OR c.product_id = ?`,
		uid, body.UserID, body.ProductID, body.ProductID).Scan(&convID)
	if err != nil {
		res, _ := db.DB.Exec("INSERT INTO conversations (product_id) VALUES (?)", nullInt64(body.ProductID))
		convID, _ = res.LastInsertId()
		db.DB.Exec("INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?), (?, ?)", convID, uid, convID, body.UserID)
	}
	c.JSON(200, gin.H{"id": convID, "product_id": body.ProductID})
}

func handleMessagesList(c *gin.Context) {
	cid := c.Param("id")
	uid := getUserID(c)
	var ok int
	if db.DB.QueryRow("SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?", cid, uid).Scan(&ok) != nil {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}
	rows, _ := db.DB.Query("SELECT m.id, m.sender_id, m.body, m.read_at, m.created_at, u.name FROM messages m JOIN users u ON u.id = m.sender_id WHERE m.conversation_id = ? ORDER BY m.created_at ASC", cid)
	var list []gin.H
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var id, senderID, created int64
			var body, name string
			var readAt sql.NullInt64
			rows.Scan(&id, &senderID, &body, &readAt, &created, &name)
			list = append(list, gin.H{"id": id, "sender_id": senderID, "body": body, "read_at": readAt.Int64, "created_at": created, "sender_name": name})
		}
	}
	c.JSON(200, list)
}

func handleMessageSend(c *gin.Context) {
	cid := c.Param("id")
	var body struct {
		Body string `json:"body"`
	}
	if c.ShouldBindJSON(&body) != nil || strings.TrimSpace(body.Body) == "" {
		c.JSON(400, gin.H{"error": "Body required"})
		return
	}
	bodyTrim := strings.TrimSpace(body.Body)
	if len(bodyTrim) > 10000 {
		bodyTrim = bodyTrim[:10000]
	}
	uid := getUserID(c)
	var ok int
	if db.DB.QueryRow("SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?", cid, uid).Scan(&ok) != nil {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}
	res, err := db.DB.Exec("INSERT INTO messages (conversation_id, sender_id, body) VALUES (?, ?, ?)", cid, uid, bodyTrim)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	db.DB.Exec("UPDATE conversations SET updated_at = unixepoch() WHERE id = ?", cid)
	mid, _ := res.LastInsertId()
	c.JSON(201, gin.H{"id": mid, "conversation_id": cid, "sender_id": uid, "body": bodyTrim, "created_at": time.Now().Unix()})
}

func handleMessageRead(c *gin.Context) {
	mid := c.Param("id")
	uid := getUserID(c)
	var senderID, convID int64
	if db.DB.QueryRow("SELECT sender_id, conversation_id FROM messages WHERE id = ?", mid).Scan(&senderID, &convID) != nil {
		c.JSON(404, gin.H{"error": "Message not found"})
		return
	}
	if senderID == uid {
		c.JSON(200, gin.H{"ok": true})
		return
	}
	var ok int
	if db.DB.QueryRow("SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?", convID, uid).Scan(&ok) != nil {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}
	db.DB.Exec("UPDATE messages SET read_at = unixepoch() WHERE id = ? AND read_at IS NULL", mid)
	c.JSON(200, gin.H{"ok": true})
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullInt64(n int64) interface{} {
	if n == 0 {
		return nil
	}
	return n
}

func mustRows(r sql.Result) int64 {
	n, _ := r.RowsAffected()
	return n
}
