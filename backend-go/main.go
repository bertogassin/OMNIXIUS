// OMNIXIUS API — Go backend. Stack: Go only (per project policy).
package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	r.Use(securityHeaders())
	r.Use(requestLogger())
	r.Use(corsMiddleware())
	r.Use(rateLimitMiddleware())
	r.Static("/uploads", cfg.UploadDir)

	r.GET("/health", func(c *gin.Context) {
		if err := db.DB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "db"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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
	auth.DELETE("/users/me", handleUserDelete)
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

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
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
		latencyMs := time.Since(start).Milliseconds()
		level := "info"
		if status >= 500 {
			level = "error"
		} else if status >= 400 {
			level = "warn"
		}
		entry := map[string]interface{}{
			"level": level, "method": method, "path": path,
			"status": status, "ip": clientIP, "latency_ms": latencyMs,
		}
		if b, err := json.Marshal(entry); err == nil {
			log.Println(string(b))
		} else {
			log.Printf("[%s] %d %s %s", method, status, path, clientIP)
		}
	}
}

func corsMiddleware() gin.HandlerFunc {
	origins := strings.Split(cfg.AllowedOrigins, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if cfg.AllowedOrigins == "" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			for _, o := range origins {
				if o != "" && (origin == o || o == "*") {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
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
	if len(body.Email) > 255 {
		c.JSON(400, gin.H{"error": "Email too long"})
		return
	}
	if len(body.Name) > 200 {
		body.Name = body.Name[:200]
	}
	user, token, err := AuthRegister(body.Email, body.Password, body.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailExists):
			c.JSON(409, gin.H{"error": "Email already registered"})
		case errors.Is(err, ErrRegistrationFailed):
			c.JSON(500, gin.H{"error": "Registration failed"})
		default:
			c.JSON(500, gin.H{"error": "Registration failed"})
		}
		return
	}
	c.JSON(201, gin.H{"user": user, "token": token})
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
	user, token, err := AuthLogin(body.Email, body.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(500, gin.H{"error": "Login failed"})
		return
	}
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

func handleUserDelete(c *gin.Context) {
	uid := getUserID(c)
	// Orders reference users; delete them first so we can delete the user.
	db.DB.Exec("DELETE FROM orders WHERE buyer_id = ? OR seller_id = ?", uid, uid)
	res, err := db.DB.Exec("DELETE FROM users WHERE id = ?", uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Account deletion failed"})
		return
	}
	if mustRows(res) == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
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
	h, err := ProductGet(c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to load product"})
		return
	}
	c.JSON(200, h)
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
	h, err := ProductCreate(getUserID(c), title, desc, category, location, imagePath, price)
	if err != nil {
		c.JSON(500, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(201, h)
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
	c.JSON(200, OrdersMy(getUserID(c)))
}

func handleOrderCreate(c *gin.Context) {
	if getUserID(c) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var body struct {
		ProductID int64 `json:"product_id"`
	}
	if c.ShouldBindJSON(&body) != nil || body.ProductID == 0 {
		c.JSON(400, gin.H{"error": "product_id required"})
		return
	}
	h, err := OrderCreate(getUserID(c), body.ProductID)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderProductNotFound):
			c.JSON(404, gin.H{"error": "Product not found"})
		case errors.Is(err, ErrOrderOwnProduct):
			c.JSON(400, gin.H{"error": "Cannot order own product"})
		default:
			c.JSON(500, gin.H{"error": "Failed"})
		}
		return
	}
	c.JSON(201, h)
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
	c.JSON(200, ConversationsList(getUserID(c)))
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
	convID, err := ConversationCreate(getUserID(c), body.UserID, body.ProductID)
	if err != nil {
		if errors.Is(err, ErrConvUserNotFound) {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(200, gin.H{"id": convID, "product_id": body.ProductID})
}

func handleMessagesList(c *gin.Context) {
	cid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid conversation id"})
		return
	}
	list, err := MessagesList(cid, getUserID(c))
	if err != nil {
		if errors.Is(err, ErrConvForbidden) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(200, list)
}

func handleMessageSend(c *gin.Context) {
	cid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid conversation id"})
		return
	}
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
	h, err := MessageSend(cid, getUserID(c), bodyTrim)
	if err != nil {
		if errors.Is(err, ErrConvForbidden) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(201, h)
}

func handleMessageRead(c *gin.Context) {
	mid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid message id"})
		return
	}
	err = MessageMarkRead(mid, getUserID(c))
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			c.JSON(404, gin.H{"error": "Message not found"})
			return
		}
		if errors.Is(err, ErrConvForbidden) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
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
