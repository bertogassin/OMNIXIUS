// OMNIXIUS API — Go backend. Stack order: Rust first, then Go (see start.bat, STACK.md).
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"omnixius-api/db"
	"omnixius-api/pqc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	if !filepath.IsAbs(cfg.SiteRoot) {
		cfg.SiteRoot, _ = filepath.Abs(cfg.SiteRoot)
	}
	cfg.SiteRoot = filepath.Clean(cfg.SiteRoot)
	log.Printf("Site root: %s", cfg.SiteRoot)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("db", "omnixius.db")
	}
	if err := db.Open(dbPath); err != nil {
		panic(err)
	}
	db.InitUploadDirs(cfg.UploadDir)
	if err := initWebAuthn(); err != nil {
		log.Printf("WebAuthn init skipped: %v (Passkeys endpoints will return 503)", err)
	}

	initWSHub()
	// Stack order: Rust first. Ping Rust service if configured.
	if cfg.RustServiceURL != "" {
		client := &http.Client{Timeout: 2 * time.Second}
		if resp, err := client.Get(cfg.RustServiceURL + "/health"); err != nil {
			log.Printf("Rust service (%s): not available (stack 1); run start-rust.bat first", cfg.RustServiceURL)
		} else {
			resp.Body.Close()
			log.Printf("Rust service (%s): ok (stack order Rust → Go)", cfg.RustServiceURL)
		}
	}

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

	r.GET("/register", handleRegisterPage)
	r.POST("/register", handleRegisterForm)
	r.GET("/login", handleLoginPage)
	r.POST("/login", handleLoginForm)

	api := r.Group("/api")
	api.POST("/auth/register", handleRegister)
	api.POST("/auth/login", handleLogin)
	api.POST("/auth/register/begin", handlePasskeyRegisterBegin)
	api.POST("/auth/register/complete", handlePasskeyRegisterComplete)
	api.POST("/auth/login/begin", handlePasskeyLoginBegin)
	api.POST("/auth/login/complete", handlePasskeyLoginComplete)
	api.GET("/auth/confirm-email", handleConfirmEmail)
	api.POST("/auth/forgot-password", handleForgotPassword)
	api.POST("/auth/reset-password", handleResetPassword)
	api.POST("/auth/recovery/verify", handleRecoveryVerify)
	api.POST("/auth/recovery/restore", handleRecoveryRestore)
	api.POST("/seed-test-user", handleSeedTestUser)

	auth := api.Group("")
	auth.Use(authRequired())
	auth.GET("/users/me", handleUserMe)
	auth.PATCH("/users/me", handleUserUpdate)
	auth.DELETE("/users/me", handleUserDelete)
	auth.GET("/auth/sessions", handleAuthSessionsList)
	auth.DELETE("/auth/sessions/:id", handleAuthSessionDelete)
	auth.GET("/auth/devices", handleAuthDevicesList)
	auth.DELETE("/auth/devices/:id", handleAuthDeviceDelete)
	auth.POST("/auth/recovery/generate", handleRecoveryGenerate)
	auth.POST("/auth/change-password", handleChangePassword)
	api.GET("/ws", handleWSWithQueryToken)
	auth.GET("/users/me/orders", handleUserOrders)
	auth.GET("/users/me/balance", handleBalanceGet)
	auth.POST("/users/me/balance/credit", handleBalanceCredit)
	auth.POST("/users/me/avatar", handleUserAvatar)
	auth.POST("/users/me/heartbeat", handleUserHeartbeat)

	api.GET("/professionals/search", handleProfessionalsSearch)

	api.GET("/products", handleProductsList)
	api.GET("/products/categories", handleProductsCategories)
	api.GET("/products/:id", handleProductGet)
	auth.GET("/products/:id/closed-content", handleProductClosedContent)
	auth.POST("/products", handleProductCreate)
	auth.PATCH("/products/:id", handleProductUpdate)
	auth.DELETE("/products/:id", handleProductDelete)
	api.GET("/products/:id/slots", handleSlotsList)
	auth.POST("/products/:id/slots", handleSlotsAdd)
	auth.POST("/products/:id/slots/:sid/book", handleSlotBook)

	api.GET("/users/:id", handleUserPublic)
	auth.POST("/subscriptions", handleSubscriptionCreate)
	auth.GET("/subscriptions/my", handleSubscriptionsMy)

	auth.GET("/orders/my", handleOrdersMy)
	auth.GET("/orders/:id", handleOrderGet)
	auth.POST("/orders", handleOrderCreate)
	auth.PATCH("/orders/:id", handleOrderUpdate)

	auth.GET("/remittances/my", handleRemittancesMy)
	auth.POST("/remittances", handleRemittanceCreate)

	auth.GET("/conversations", handleConversationsList)
	auth.GET("/conversations/unread-count", handleConversationsUnreadCount)
	auth.GET("/conversations/:id", handleConversationGet)
	auth.POST("/conversations", handleConversationCreate)
	auth.GET("/messages/conversation/:id", handleMessagesList)
	auth.POST("/messages/conversation/:id", handleMessageSend)
	auth.POST("/messages/:id/read", handleMessageRead)

	// Vault API v1 (ARCHITECTURE-V4)
	vault := api.Group("/v1/vault", authRequired())
	vault.POST("/files/upload-url", handleVaultUploadURL)       // 501, for future S3 pre-signed
	vault.POST("/files/:id/complete", handleVaultCompleteUpload) // 501
	vault.GET("/files/:id/download-url", handleVaultDownloadURL) // 501
	vault.POST("/files", handleVaultUploadFile)
	vault.GET("/files", handleVaultListFiles)
	vault.GET("/files/:id", handleVaultGetFile)
	vault.GET("/files/:id/download", handleVaultDownloadFile)
	vault.DELETE("/files/:id", handleVaultDeleteFile)
	vault.POST("/folders", handleVaultCreateFolder)
	vault.GET("/folders", handleVaultListFolders)
	vault.DELETE("/folders/:id", handleVaultDeleteFolder)
	vault.POST("/search", handleVaultSearch)

	// Wallet (§15)
	auth.GET("/wallet/balances", handleWalletBalances)
	auth.GET("/wallet/balances/:currency", handleWalletBalanceByCurrency)
	auth.GET("/wallet/transactions", handleWalletTransactions)
	auth.GET("/wallet/transactions/:id", handleWalletTransactionByID)
	auth.POST("/wallet/transfer", handleWalletTransfer)
	auth.POST("/wallet/transfer/verify", handleWalletTransferVerify)
	auth.GET("/wallet/deposit/addresses", handleWalletDepositAddressesList)
	auth.POST("/wallet/deposit/addresses", handleWalletDepositAddressCreate)
	auth.POST("/wallet/hold", handleWalletHold)
	auth.POST("/wallet/hold/:id/release", handleWalletHoldRelease)
	auth.POST("/wallet/hold/:id/capture", handleWalletHoldCapture)
	auth.POST("/wallet/export", handleWalletExport)
	auth.POST("/wallet/import", handleWalletImport)

	// Notifications (§16)
	auth.GET("/notifications/settings", handleNotificationsSettingsGet)
	auth.PATCH("/notifications/settings", handleNotificationsSettingsPatch)
	auth.GET("/notifications/history", handleNotificationsHistory)
	auth.GET("/notifications/history/:id", handleNotificationsHistoryByID)
	auth.POST("/notifications/history/:id/read", handleNotificationsHistoryRead)
	auth.POST("/notifications/push/tokens", handleNotificationsPushTokenCreate)
	auth.DELETE("/notifications/push/tokens/:id", handleNotificationsPushTokenDelete)
	auth.POST("/notifications/test", handleNotificationsTest)

	// Admin (§18) — requires admin role
	adminGroup := api.Group("/admin", authRequired(), adminRequired())
	adminGroup.GET("/stats", handleAdminStats)
	adminGroup.GET("/reports", handleAdminReportsList)
	adminGroup.GET("/reports/:id", handleAdminReportGet)
	adminGroup.POST("/reports/:id/assign", handleAdminReportAssign)
	adminGroup.POST("/reports/:id/resolve", handleAdminReportResolve)
	adminGroup.GET("/users/:id", handleAdminUserGet)
	adminGroup.POST("/users/:id/ban", handleAdminUserBan)
	adminGroup.POST("/users/:id/unban", handleAdminUserUnban)
	api.POST("/reports", authRequired(), handleReportCreate)

	spaRoot := filepath.Join(cfg.SiteRoot, "web", "dist")
	r.NoRoute(staticSiteHandler(cfg.SiteRoot, spaRoot))

	port := ":" + cfg.Port
	if err := r.Run(port); err != nil {
		panic(err)
	}
}

func staticSiteHandler(siteRoot, spaRoot string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/uploads") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		// SPA (TypeScript app) at /app — serve from web/dist, fallback to index.html
		if strings.HasPrefix(path, "/app") {
			rel := strings.TrimPrefix(path, "/app")
			if rel == "" {
				rel = "/"
			}
			rel = strings.TrimPrefix(rel, "/")
			if rel == "" {
				rel = "index.html"
			}
			cleanRel := filepath.Clean(filepath.FromSlash(rel))
			if cleanRel == "" || cleanRel == "." || strings.HasPrefix(cleanRel, "..") {
				cleanRel = "index.html"
			}
			fullPath := filepath.Join(spaRoot, cleanRel)
			if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
				c.File(fullPath)
				return
			}
			indexPath := filepath.Join(spaRoot, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
				return
			}
		}
		if path == "/" {
			path = "index.html"
		} else {
			path = strings.TrimPrefix(path, "/")
		}
		cleanPath := filepath.Clean(filepath.FromSlash(path))
		if cleanPath == "" || cleanPath == "." || strings.HasPrefix(cleanPath, "..") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		fullPath := filepath.Join(siteRoot, cleanPath)
		absRoot, _ := filepath.Abs(siteRoot)
		absFull, err := filepath.Abs(fullPath)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		prefix := absRoot + string(filepath.Separator)
		if absFull != absRoot && !strings.HasPrefix(absFull, prefix) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.File(fullPath)
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
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
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
			"request_id": requestID, "level": level, "method": method, "path": path,
			"status": status, "ip": clientIP, "duration_ms": latencyMs,
		}
		if uid, ok := c.Get("userID"); ok && uid != nil {
			if id, ok := uid.(int64); ok {
				entry["user_id"] = id
			}
		}
		if b, err := json.Marshal(entry); err == nil {
			log.Println(string(b))
		} else {
			log.Printf("[%s] %s %d %s %s", requestID, method, status, path, clientIP)
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
		uid, _, sessionID, err := pqc.VerifyToken(cfg.PQCPublicKey, tok)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if sessionID != 0 {
			var n int
			if db.DB.QueryRow("SELECT 1 FROM sessions WHERE id = ? AND user_id = ? AND expires_at > ?", sessionID, uid, time.Now().Unix()).Scan(&n) != nil || n == 0 {
				c.JSON(401, gin.H{"error": "Session invalid or expired"})
				c.Abort()
				return
			}
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
		c.Set("userRole", role)
		c.Set("userName", name)
		c.Set("userAvatar", avatar)
		_, _ = db.DB.Exec("UPDATE users SET last_seen_at = unixepoch() WHERE id = ?", uid)
		c.Next()
	}
}

// handleWSWithQueryToken validates token from query (for WebSocket, browser cannot send Authorization header) then runs handleWS.
func handleWSWithQueryToken(c *gin.Context) {
	tok := strings.TrimSpace(c.Query("token"))
	if tok == "" {
		c.JSON(401, gin.H{"error": "token required"})
		return
	}
	uid, _, sessionID, err := pqc.VerifyToken(cfg.PQCPublicKey, tok)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}
	if sessionID != 0 {
		var n int
		if db.DB.QueryRow("SELECT 1 FROM sessions WHERE id = ? AND user_id = ? AND expires_at > ?", sessionID, uid, time.Now().Unix()).Scan(&n) != nil || n == 0 {
			c.JSON(401, gin.H{"error": "Session invalid or expired"})
			return
		}
	}
	c.Set("userID", uid)
	handleWS(c)
}

func adminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, _ := c.Get("userRole")
		if role, ok := r.(string); !ok || role != "admin" {
			c.JSON(403, gin.H{"error": "Admin required"})
			c.Abort()
			return
		}
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

func auditLog(userID int64, action, entityType, entityID, details string) {
	_, _ = db.DB.Exec("INSERT INTO audit_log (user_id, action, entity_type, entity_id, details) VALUES (?, ?, ?, ?, ?)",
		userID, action, entityType, entityID, details)
}

// getOptionalUserID returns user ID if valid Bearer token present; otherwise 0 (for optional-auth routes).
func getOptionalUserID(c *gin.Context) int64 {
	if v, _ := c.Get("userID"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	tok := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tok == "" {
		return 0
	}
	uid, _, sessionID, err := pqc.VerifyToken(cfg.PQCPublicKey, tok)
	if err != nil {
		return 0
	}
	if sessionID != 0 {
		var n int
		if db.DB.QueryRow("SELECT 1 FROM sessions WHERE id = ? AND user_id = ? AND expires_at > ?", sessionID, uid, time.Now().Unix()).Scan(&n) != nil || n == 0 {
			return 0
		}
	}
	var n int
	if db.DB.QueryRow("SELECT 1 FROM users WHERE id = ?", uid).Scan(&n) != nil {
		return 0
	}
	return uid
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

const (
	maxEmailLen    = 255
	maxPasswordLen = 128
	minPasswordLen = 8
	maxNameLen     = 200
)

func isValidEmail(s string) bool {
	if s == "" || len(s) > maxEmailLen {
		return false
	}
	if strings.Contains(s, " ") {
		return false
	}
	at := strings.LastIndex(s, "@")
	if at <= 0 || at == len(s)-1 {
		return false
	}
	domain := s[at+1:]
	if !strings.Contains(domain, ".") || len(domain) < 2 {
		return false
	}
	return true
}

func handleRegister(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Email and password required"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if body.Email == "" {
		c.JSON(400, gin.H{"error": "Email required"})
		return
	}
	if len(body.Email) > maxEmailLen {
		c.JSON(400, gin.H{"error": "Email must be up to 255 characters"})
		return
	}
	if !isValidEmail(body.Email) {
		c.JSON(400, gin.H{"error": "Invalid email format. Use: letters, numbers, @ and a dot (e.g. name@domain.com)"})
		return
	}
	if len(body.Password) < minPasswordLen {
		c.JSON(400, gin.H{"error": "Password must be at least 8 characters"})
		return
	}
	if len(body.Password) > maxPasswordLen {
		c.JSON(400, gin.H{"error": "Password must be up to 128 characters"})
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if len(body.Name) > maxNameLen {
		body.Name = body.Name[:maxNameLen]
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

// handleSeedTestUser creates a test user only when the DB has zero users. For local dev.
const testUserEmail = "test@test.com"
const testUserPassword = "Test123!"

func handleSeedTestUser(c *gin.Context) {
	var n int
	if db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&n) != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	if n > 0 {
		c.JSON(200, gin.H{"ok": true, "message": "Users already exist. Use existing account or register.", "test_email": testUserEmail, "test_password": testUserPassword})
		return
	}
	user, token, err := AuthRegister(testUserEmail, testUserPassword, "Test User")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"ok": true, "message": "Test user created. You can sign in.", "user": user, "token": token, "test_email": testUserEmail, "test_password": testUserPassword})
}

	const registerPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Register — OMNIXIUS</title>
  <style>
    body { font-family: system-ui, sans-serif; max-width: 420px; margin: 2rem auto; padding: 1rem; background: #0c0c0f; color: #e8e6e3; }
    h1 { font-size: 1.5rem; margin-bottom: 0.5rem; }
    p.sub { color: #8a8a8a; font-size: 0.95rem; margin-bottom: 1.5rem; }
    label { display: block; margin-bottom: 0.25rem; color: #8a8a8a; }
    input { width: 100%; padding: 0.6rem; margin-bottom: 1rem; background: #14141a; border: 1px solid #2a2a32; border-radius: 6px; color: #e8e6e3; box-sizing: border-box; }
    .err { color: #e74c3c; font-size: 0.9rem; margin-bottom: 0.5rem; }
    button { width: 100%; padding: 0.75rem; background: #00d4aa; color: #0c0c0f; border: none; border-radius: 8px; font-weight: 600; cursor: pointer; font-size: 1rem; }
    button:hover { opacity: 0.9; }
    a { color: #00d4aa; }
    .link { text-align: center; margin-top: 1rem; }
  </style>
</head>
<body>
  <h1>Register</h1>
  <p class="sub">Create your OMNIXIUS account (via Go backend).</p>
  {{.ErrorHTML}}
  <form method="POST" action="/register">
    <label>Email</label>
    <input type="email" name="email" required value="{{.Email}}" autocomplete="email">
    <label>Password (min 8 characters)</label>
    <input type="password" name="password" required minlength="8" autocomplete="new-password">
    <label>Confirm password</label>
    <input type="password" name="password2" required minlength="8" autocomplete="new-password">
    <label>Name (optional)</label>
    <input type="text" name="name" value="{{.Name}}" autocomplete="name" maxlength="200">
    <button type="submit">Create account</button>
  </form>
  <p class="link"><a href="/login">Already have an account? Sign in</a></p>
</body>
</html>`

const loginPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Sign in — OMNIXIUS</title>
  <style>
    body { font-family: system-ui, sans-serif; max-width: 420px; margin: 2rem auto; padding: 1rem; background: #0c0c0f; color: #e8e6e3; }
    h1 { font-size: 1.5rem; margin-bottom: 1.5rem; }
    label { display: block; margin-bottom: 0.25rem; color: #8a8a8a; }
    input { width: 100%; padding: 0.6rem; margin-bottom: 1rem; background: #14141a; border: 1px solid #2a2a32; border-radius: 6px; color: #e8e6e3; box-sizing: border-box; }
    .err { color: #e74c3c; font-size: 0.9rem; margin-bottom: 0.5rem; }
    button { width: 100%; padding: 0.75rem; background: #00d4aa; color: #0c0c0f; border: none; border-radius: 8px; font-weight: 600; cursor: pointer; font-size: 1rem; }
    a { color: #00d4aa; }
    .link { text-align: center; margin-top: 1rem; }
  </style>
</head>
<body>
  <h1>Sign in</h1>
  {{.Error}}
  <form method="POST" action="/login">
    <label>Email</label>
    <input type="email" name="email" required value="{{.Email}}" autocomplete="email">
    <label>Password</label>
    <input type="password" name="password" required autocomplete="current-password">
    <button type="submit">Sign in</button>
  </form>
  <p class="link"><a href="/register">No account? Register</a></p>
</body>
</html>`

func handleRegisterPage(c *gin.Context) {
	html := strings.ReplaceAll(registerPageHTML, "{{.ErrorHTML}}", "")
	html = strings.ReplaceAll(html, "{{.Email}}", "")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, strings.ReplaceAll(html, "{{.Name}}", ""))
}

func handleRegisterForm(c *gin.Context) {
	email := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	password := c.PostForm("password")
	password2 := c.PostForm("password2")
	name := strings.TrimSpace(c.PostForm("name"))
	if len(name) > maxNameLen {
		name = name[:maxNameLen]
	}

	if email == "" {
		serveRegisterError(c, "Email required", email, name)
		return
	}
	if !isValidEmail(email) {
		serveRegisterError(c, "Invalid email format", email, name)
		return
	}
	if len(password) < minPasswordLen {
		serveRegisterError(c, "Password must be at least 8 characters", email, name)
		return
	}
	if len(password) > maxPasswordLen {
		serveRegisterError(c, "Password must be up to 128 characters", email, name)
		return
	}
	if password != password2 {
		serveRegisterError(c, "Passwords do not match", email, name)
		return
	}

	user, token, err := AuthRegister(email, password, name)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			serveRegisterError(c, "Email already registered", email, name)
			return
		}
		serveRegisterError(c, "Registration failed. Try again.", email, name)
		return
	}

	if cfg.AppURL != "" {
		apiBase := c.Request.URL.Scheme + "://" + c.Request.Host
		redirectURL := cfg.AppURL + "/app/dashboard.html?token=" + url.QueryEscape(token) + "&api_url=" + url.QueryEscape(apiBase)
		if u, ok := user["name"]; ok && u != nil {
			redirectURL += "&name=" + url.QueryEscape(stringOrEmpty(user["name"]))
		}
		if e, ok := user["email"]; ok && e != nil {
			redirectURL += "&email=" + url.QueryEscape(stringOrEmpty(user["email"]))
		}
		c.Redirect(http.StatusFound, redirectURL)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, `<html><body><p>Account created. Token: `+token+`</p><p><a href="/register">Back to register</a></p></body></html>`)
}

func stringOrEmpty(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func serveRegisterError(c *gin.Context, errMsg, email, name string) {
	errBlock := "<p class=\"err\">" + templateHTMLEscape(errMsg) + "</p>"
	html := strings.ReplaceAll(registerPageHTML, "{{.ErrorHTML}}", errBlock)
	html = strings.ReplaceAll(html, "{{.Email}}", templateHTMLEscape(email))
	html = strings.ReplaceAll(html, "{{.Name}}", templateHTMLEscape(name))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(400, html)
}

func templateHTMLEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

func handleLoginPage(c *gin.Context) {
	html := strings.ReplaceAll(loginPageHTML, "{{.Error}}", "")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, strings.ReplaceAll(html, "{{.Email}}", ""))
}

func handleLoginForm(c *gin.Context) {
	ip := c.ClientIP()
	limiter := getLoginLimiter(ip)
	if !limiter.Allow() {
		serveLoginError(c, "Too many attempts. Try again later.", c.PostForm("email"))
		return
	}
	email := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	password := c.PostForm("password")
	if email == "" {
		serveLoginError(c, "Email required", "")
		return
	}
	user, token, err := AuthLogin(email, password)
	if err != nil {
		serveLoginError(c, "Invalid email or password", email)
		return
	}
	if cfg.AppURL != "" {
		apiBase := c.Request.URL.Scheme + "://" + c.Request.Host
		redirectURL := cfg.AppURL + "/app/dashboard.html?token=" + url.QueryEscape(token) + "&api_url=" + url.QueryEscape(apiBase)
		if u, ok := user["name"]; ok && u != nil {
			redirectURL += "&name=" + url.QueryEscape(stringOrEmpty(user["name"]))
		}
		redirectURL += "&email=" + url.QueryEscape(stringOrEmpty(user["email"]))
		c.Redirect(http.StatusFound, redirectURL)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, `<html><body><p>Signed in. Token: `+token+`</p><p><a href="/login">Back to login</a></p></body></html>`)
}

func serveLoginError(c *gin.Context, errMsg, email string) {
	errBlock := "<p class=\"err\">" + templateHTMLEscape(errMsg) + "</p>"
	html := strings.ReplaceAll(loginPageHTML, "{{.Error}}", errBlock)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(401, strings.ReplaceAll(html, "{{.Email}}", templateHTMLEscape(email)))
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
	if body.Email == "" {
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
	var emailVerified, phoneVerified int
	if db.DB.QueryRow("SELECT email, role, name, avatar_path, COALESCE(email_verified, 0), COALESCE(phone_verified, 0) FROM users WHERE id = ?", id).Scan(&email, &role, &name, &avatar, &emailVerified, &phoneVerified) != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	verified := emailVerified == 1 || phoneVerified == 1
	c.JSON(200, gin.H{"id": id, "email": email, "role": role, "name": name, "avatar_path": avatar.String, "email_verified": emailVerified == 1, "phone_verified": phoneVerified == 1, "verified": verified})
}

func handleUserUpdate(c *gin.Context) {
	var body struct {
		Name         *string  `json:"name"`
		ProfessionID *string  `json:"profession_id"`
		Lat          *float64 `json:"lat"`
		Lng          *float64 `json:"lng"`
	}
	c.ShouldBindJSON(&body)
	uid := getUserID(c)
	if body.Name != nil {
		name := strings.TrimSpace(*body.Name)
		if len(name) > 200 {
			name = name[:200]
		}
		db.DB.Exec("UPDATE users SET name = ?, updated_at = unixepoch() WHERE id = ?", name, uid)
	}
	if body.ProfessionID != nil {
		p := strings.TrimSpace(*body.ProfessionID)
		if len(p) > 64 {
			p = p[:64]
		}
		db.DB.Exec("UPDATE users SET profession_id = ?, updated_at = unixepoch() WHERE id = ?", nullStr(p), uid)
	}
	if body.Lat != nil && body.Lng != nil {
		db.DB.Exec("UPDATE users SET lat = ?, lng = ?, updated_at = unixepoch() WHERE id = ?", *body.Lat, *body.Lng, uid)
	}
	handleUserMe(c)
}

func handleUserHeartbeat(c *gin.Context) {
	uid := getUserID(c)
	db.DB.Exec("UPDATE users SET last_seen_at = unixepoch() WHERE id = ?", uid)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// haversineKm returns distance in km between two points (lat/lng in degrees).
func haversineKm(lat1, lng1, lat2, lng2 float64) float64 {
	const earthR = 6371 // km
	rad := func(d float64) float64 { return d * (3.141592653589793 / 180) }
	dLat := rad(lat2 - lat1)
	dLng := rad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(rad(lat1))*math.Cos(rad(lat2))*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthR * c
}

func handleProfessionalsSearch(c *gin.Context) {
	profession := strings.TrimSpace(c.Query("profession"))
	onlineOnly := c.Query("online") == "1" || c.Query("online") == "true"
	latQ, lngQ := c.Query("lat"), c.Query("lng")
	var latCenter, lngCenter float64
	useLocation := false
	if latQ != "" && lngQ != "" {
		if la, err := strconv.ParseFloat(latQ, 64); err == nil && la >= -90 && la <= 90 {
			if ln, err := strconv.ParseFloat(lngQ, 64); err == nil && ln >= -180 && ln <= 180 {
				latCenter, lngCenter = la, ln
				useLocation = true
			}
		}
	}
	radiusKm := 0.0
	if r := c.Query("radius_km"); r != "" {
		if f, err := strconv.ParseFloat(r, 64); err == nil && f > 0 {
			radiusKm = f
		}
	}
	sortBy := strings.TrimSpace(c.Query("sort")) // "rating" | "distance" | ""
	limit := 50
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	now := time.Now().Unix()
	onlineSince := now - 900
	sel := `SELECT id, name, avatar_path, COALESCE(profession_id,''), lat, lng, last_seen_at, COALESCE(rating_avg,0), COALESCE(rating_count,0) FROM users WHERE profession_id IS NOT NULL AND profession_id != ''`
	args := []interface{}{}
	if profession != "" {
		sel += ` AND profession_id = ?`
		args = append(args, profession)
	}
	if onlineOnly {
		sel += ` AND last_seen_at > ?`
		args = append(args, onlineSince)
	}
	sel += ` ORDER BY COALESCE(last_seen_at,0) DESC`
	args = append(args, limit*3) // fetch extra for distance filter
	sel += ` LIMIT ?`
	rows, err := db.DB.Query(sel, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"professionals": []gin.H{}})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id int64
		var name, avatarPath, professionID string
		var lat, lng sql.NullFloat64
		var lastSeenAt sql.NullInt64
		var ratingAvg float64
		var ratingCount int64
		if rows.Scan(&id, &name, &avatarPath, &professionID, &lat, &lng, &lastSeenAt, &ratingAvg, &ratingCount) != nil {
			continue
		}
		ls := lastSeenAt.Int64
		item := gin.H{"id": id, "name": name, "avatar_path": avatarPath, "profession_id": professionID, "last_seen_at": ls, "online": ls > onlineSince, "rating_avg": ratingAvg, "rating_count": ratingCount}
		if lat.Valid {
			item["lat"] = lat.Float64
		}
		if lng.Valid {
			item["lng"] = lng.Float64
		}
		if useLocation && lat.Valid && lng.Valid {
			km := haversineKm(latCenter, lngCenter, lat.Float64, lng.Float64)
			item["distance_km"] = math.Round(km*10) / 10
			if radiusKm > 0 && km > radiusKm {
				continue
			}
		}
		list = append(list, item)
	}
	if list == nil {
		list = []gin.H{}
	}
	if useLocation && sortBy == "distance" {
		sort.Slice(list, func(i, j int) bool {
			di, _ := list[i]["distance_km"].(float64)
			dj, _ := list[j]["distance_km"].(float64)
			return di < dj
		})
	} else if sortBy == "rating" {
		sort.Slice(list, func(i, j int) bool {
			ri, _ := list[i]["rating_avg"].(float64)
			rj, _ := list[j]["rating_avg"].(float64)
			return ri > rj
		})
	}
	// Stack: Rust first. When Rust is up, use it for ranking (heavy compute path).
	if cfg.RustServiceURL != "" && sortBy != "" && len(list) > 0 {
		payload := map[string]interface{}{"items": list, "sort": sortBy}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, cfg.RustServiceURL+"/rank", strings.NewReader(string(body)))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{Timeout: 3 * time.Second}
			if resp, err := client.Do(req); err == nil && resp.StatusCode == http.StatusOK {
				var out struct {
					Items []gin.H `json:"items"`
				}
				if json.NewDecoder(resp.Body).Decode(&out) == nil && len(out.Items) > 0 {
					list = out.Items
				}
				resp.Body.Close()
			}
		}
	}
	if len(list) > limit {
		list = list[:limit]
	}
	c.JSON(http.StatusOK, gin.H{"professionals": list})
}

func handleChangePassword(c *gin.Context) {
	if getUserID(c) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.CurrentPassword == "" || body.NewPassword == "" {
		c.JSON(400, gin.H{"error": "current_password and new_password required"})
		return
	}
	if len(body.NewPassword) < minPasswordLen {
		c.JSON(400, gin.H{"error": "New password must be at least 8 characters"})
		return
	}
	if len(body.NewPassword) > maxPasswordLen {
		c.JSON(400, gin.H{"error": "New password must be up to 128 characters"})
		return
	}
	uid := getUserID(c)
	var hash string
	if db.DB.QueryRow("SELECT password_hash FROM users WHERE id = ?", uid).Scan(&hash) != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if !checkPassword(hash, body.CurrentPassword) {
		c.JSON(401, gin.H{"error": "Current password is incorrect"})
		return
	}
	newHash := hashPasswordArgon2(body.NewPassword)
	db.DB.Exec("UPDATE users SET password_hash = ?, updated_at = unixepoch() WHERE id = ?", newHash, uid)
	c.JSON(200, gin.H{"ok": true})
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

func handleBalanceGet(c *gin.Context) {
	c.JSON(200, BalanceGet(getUserID(c)))
}

func handleBalanceCredit(c *gin.Context) {
	var body struct {
		Amount float64 `json:"amount"`
	}
	if c.ShouldBindJSON(&body) != nil || body.Amount <= 0 {
		c.JSON(400, gin.H{"error": "amount required (positive number)"})
		return
	}
	if body.Amount > 1e9 {
		c.JSON(400, gin.H{"error": "amount too large"})
		return
	}
	h, err := BalanceCredit(getUserID(c), body.Amount)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(200, h)
}

func handleUserOrders(c *gin.Context) {
	id := getUserID(c)
	// asBuyer
	rows, _ := db.DB.Query(`SELECT o.id, o.status, o.created_at, COALESCE(o.installment_plan, ''), p.title, p.price, p.image_path, u.name FROM orders o JOIN products p ON p.id = o.product_id JOIN users u ON u.id = o.seller_id WHERE o.buyer_id = ? ORDER BY o.created_at DESC`, id)
	asBuyer := rowsToOrderList(rows)
	rows, _ = db.DB.Query(`SELECT o.id, o.status, o.created_at, COALESCE(o.installment_plan, ''), p.title, p.price, p.image_path, u.name FROM orders o JOIN products p ON p.id = o.product_id JOIN users u ON u.id = o.buyer_id WHERE o.seller_id = ? ORDER BY o.created_at DESC`, id)
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
		var status, installmentPlan, title string
		var price float64
		var imagePath, name sql.NullString
		rows.Scan(&id, &status, &created, &installmentPlan, &title, &price, &imagePath, &name)
		out = append(out, gin.H{"id": id, "status": status, "created_at": created, "installment_plan": installmentPlan, "title": title, "price": price, "image_path": imagePath.String, "seller_name": name.String})
	}
	return out
}

const maxAvatarBytes = 5 * 1024 * 1024   // 5 MB
const maxProductImageBytes = 10 * 1024 * 1024 // 10 MB

func handleUserAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, gin.H{"error": "File required"})
		return
	}
	if file.Size > maxAvatarBytes {
		c.JSON(400, gin.H{"error": "File too large (max 5 MB)"})
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
	qry := `SELECT p.id, p.title, p.price, p.category, p.location, p.image_path, COALESCE(p.is_service, 0), COALESCE(p.is_subscription, 0), p.created_at, u.id, u.name, COALESCE(u.email_verified, 0), COALESCE(u.phone_verified, 0) FROM products p JOIN users u ON u.id = p.user_id WHERE 1=1`
	args := []interface{}{}
	if uid := c.Query("user_id"); uid != "" {
		if uidNum, err := strconv.ParseInt(uid, 10, 64); err == nil {
			qry += ` AND p.user_id = ?`
			args = append(args, uidNum)
		}
	}
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
	if v := c.Query("service"); v == "1" || v == "true" {
		qry += ` AND COALESCE(p.is_service, 0) = 1`
	}
	if v := c.Query("subscription"); v == "1" || v == "true" {
		qry += ` AND COALESCE(p.is_subscription, 0) = 1`
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
			var isService, isSubscription int64
			var sellerID int64
			var emailVerified, phoneVerified int64
			rows.Scan(&id, &title, &price, &category, &location, &imagePath, &isService, &isSubscription, &created, &sellerID, &sellerName, &emailVerified, &phoneVerified)
			sellerVerified := emailVerified == 1 || phoneVerified == 1
			list = append(list, gin.H{"id": id, "title": title, "price": price, "category": category, "location": location, "image_path": imagePath.String, "is_service": isService, "is_subscription": isSubscription, "created_at": created, "seller_id": sellerID, "seller_name": sellerName.String, "seller_verified": sellerVerified})
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

func handleProductClosedContent(c *gin.Context) {
	idStr := c.Param("id")
	url, ownerID, isSub, err := ProductClosedContent(idStr)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	uid := getUserID(c)
	if uid == 0 {
		c.JSON(401, gin.H{"error": "Sign in to access closed content"})
		return
	}
	pid, _ := strconv.ParseInt(idStr, 10, 64)
	if ownerID == uid || IsSubscribed(pid, uid) {
		c.JSON(200, gin.H{"url": url})
		return
	}
	if isSub != 1 {
		c.JSON(400, gin.H{"error": "This listing has no closed content"})
		return
	}
	c.JSON(403, gin.H{"error": "Subscribe to access this content"})
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
		if file.Size > maxProductImageBytes {
			c.JSON(400, gin.H{"error": "Image too large (max 10 MB)"})
			return
		}
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
	isService := 0
	if v := c.PostForm("is_service"); v == "1" || v == "on" || v == "true" {
		isService = 1
	}
	isSubscription := 0
	if v := c.PostForm("is_subscription"); v == "1" || v == "on" || v == "true" {
		isSubscription = 1
	}
	closedContentURL := strings.TrimSpace(c.PostForm("closed_content_url"))
	h, err := ProductCreate(getUserID(c), title, desc, category, location, imagePath, price, isService, isSubscription, closedContentURL)
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
	if v := c.PostForm("is_service"); v != "" {
		isService := 0
		if v == "1" || v == "on" || v == "true" {
			isService = 1
		}
		db.DB.Exec("UPDATE products SET is_service = ? WHERE id = ?", isService, id)
	}
	if v := c.PostForm("is_subscription"); v != "" {
		isSub := 0
		if v == "1" || v == "on" || v == "true" {
			isSub = 1
		}
		db.DB.Exec("UPDATE products SET is_subscription = ? WHERE id = ?", isSub, id)
	}
	if closedURL, ok := c.GetPostForm("closed_content_url"); ok {
		u := strings.TrimSpace(closedURL)
		if len(u) > 2048 {
			u = u[:2048]
		}
		db.DB.Exec("UPDATE products SET closed_content_url = ? WHERE id = ?", nullStr(u), id)
	}
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

func handleSlotsList(c *gin.Context) {
	id := c.Param("id")
	list, err := SlotsList(id, getOptionalUserID(c))
	if err != nil {
		if errors.Is(err, ErrSlotProductNotFound) {
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(200, list)
}

func handleSlotsAdd(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		SlotAt int64 `json:"slot_at"`
	}
	c.ShouldBindJSON(&body)
	if body.SlotAt <= 0 {
		body.SlotAt = time.Now().Unix()
	}
	h, err := SlotsAdd(id, getUserID(c), body.SlotAt)
	if err != nil {
		if errors.Is(err, ErrSlotProductNotFound) {
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}
		if errors.Is(err, ErrSlotForbidden) {
			c.JSON(403, gin.H{"error": "Only the product owner can add slots"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(201, h)
}

func handleSlotBook(c *gin.Context) {
	productID := c.Param("id")
	slotID := c.Param("sid")
	h, err := SlotBook(productID, slotID, getUserID(c))
	if err != nil {
		if errors.Is(err, ErrSlotProductNotFound) || errors.Is(err, ErrSlotNotFound) {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		if errors.Is(err, ErrSlotNotFree) {
			c.JSON(400, gin.H{"error": "Slot already booked"})
			return
		}
		if errors.Is(err, ErrSlotNotService) {
			c.JSON(400, gin.H{"error": "Booking is only for service listings"})
			return
		}
		if errors.Is(err, ErrOrderOwnProduct) {
			c.JSON(400, gin.H{"error": "Cannot book your own service"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(201, h)
}

func handleUserPublic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	var name, avatarPath sql.NullString
	var emailVerified, phoneVerified int
	if db.DB.QueryRow("SELECT name, avatar_path, COALESCE(email_verified, 0), COALESCE(phone_verified, 0) FROM users WHERE id = ?", id).Scan(&name, &avatarPath, &emailVerified, &phoneVerified) != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	verified := emailVerified == 1 || phoneVerified == 1
	c.JSON(200, gin.H{"id": id, "name": name.String, "avatar_path": avatarPath.String, "verified": verified})
}

func handleSubscriptionCreate(c *gin.Context) {
	var body struct {
		ProductID int64 `json:"product_id"`
	}
	if c.ShouldBindJSON(&body) != nil || body.ProductID == 0 {
		c.JSON(400, gin.H{"error": "product_id required"})
		return
	}
	h, err := Subscribe(strconv.FormatInt(body.ProductID, 10), getUserID(c))
	if err != nil {
		if errors.Is(err, ErrSubProductNotFound) {
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}
		if errors.Is(err, ErrSubNotSubscription) {
			c.JSON(400, gin.H{"error": "Product is not a subscription listing"})
			return
		}
		if errors.Is(err, ErrSubOwnProduct) {
			c.JSON(400, gin.H{"error": "Cannot subscribe to your own listing"})
			return
		}
		if errors.Is(err, ErrSubAlreadySubscribed) {
			c.JSON(400, gin.H{"error": "Already subscribed"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(201, h)
}

func handleSubscriptionsMy(c *gin.Context) {
	c.JSON(200, SubscriptionsMy(getUserID(c)))
}

func handleOrdersMy(c *gin.Context) {
	c.JSON(200, OrdersMy(getUserID(c)))
}

func handleOrderGet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	uid := getUserID(c)
	var oid, pid, buyerID, sellerID, createdAt int64
	var status, title string
	var price float64
	var img sql.NullString
	var urgent int
	err = db.DB.QueryRow(
		`SELECT o.id, o.product_id, o.buyer_id, o.seller_id, o.status, o.created_at, COALESCE(o.urgent, 0), p.title, p.price, p.image_path
		 FROM orders o JOIN products p ON p.id = o.product_id WHERE o.id = ? AND (o.buyer_id = ? OR o.seller_id = ?)`,
		id, uid, uid,
	).Scan(&oid, &pid, &buyerID, &sellerID, &status, &createdAt, &urgent, &title, &price, &img)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": oid, "product_id": pid, "buyer_id": buyerID, "seller_id": sellerID,
		"status": status, "created_at": createdAt, "urgent": urgent == 1, "title": title, "price": price, "image_path": img.String,
	})
}

func handleOrderCreate(c *gin.Context) {
	if getUserID(c) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var body struct {
		ProductID       int64  `json:"product_id"`
		InstallmentPlan string `json:"installment_plan"`
		Urgent          bool   `json:"urgent"`
	}
	if c.ShouldBindJSON(&body) != nil || body.ProductID == 0 {
		c.JSON(400, gin.H{"error": "product_id required"})
		return
	}
	installmentPlan := ""
	if body.InstallmentPlan == "requested" || body.InstallmentPlan == "installments" {
		installmentPlan = "requested"
	}
	uid := getUserID(c)
	h, err := OrderCreate(uid, body.ProductID, installmentPlan, body.Urgent)
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
	// Notify seller about new order
	if sellerID, ok := h["seller_id"].(int64); ok && sellerID != 0 {
		if oid, ok2 := h["id"].(int64); ok2 && oid != 0 {
			dataJSON := `{"order_id":` + strconv.FormatInt(oid, 10) + `}`
			title := "New order"
			if body.Urgent {
				title = "Urgent order"
			}
			_, _ = db.DB.Exec(
				"INSERT INTO notifications_queue (user_id, type, channel, title, body, data, status) VALUES (?, 'order_new', 'in_app', ?, ?, ?, 'pending')",
				sellerID, title, "Order #"+strconv.FormatInt(oid, 10)+" — view details and accept or decline.", dataJSON,
			)
			BroadcastToUser(sellerID, "notification", gin.H{"type": "order_new", "title": title, "order_id": oid})
		}
	}
	auditLog(uid, "order_created", "order", strconv.FormatInt(h["id"].(int64), 10), "")
	c.JSON(201, h)
}

func handleOrderUpdate(c *gin.Context) {
	idStr := c.Param("id")
	var body struct {
		Status          string `json:"status"`
		InstallmentPlan string `json:"installment_plan"`
	}
	c.ShouldBindJSON(&body)
	var buyer, seller int64
	var currentStatus string
	if db.DB.QueryRow("SELECT buyer_id, seller_id, status FROM orders WHERE id = ?", idStr).Scan(&buyer, &seller, &currentStatus) != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}
	uid := getUserID(c)
	if buyer != uid && seller != uid {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}
	if body.Status != "" {
		if body.Status != "pending" && body.Status != "confirmed" && body.Status != "completed" && body.Status != "cancelled" {
			c.JSON(400, gin.H{"error": "Invalid status"})
			return
		}
		// State machine: pending -> confirmed|cancelled; confirmed -> completed; cancelled/completed terminal
		allowed := false
		switch currentStatus {
		case "pending":
			allowed = body.Status == "confirmed" || body.Status == "cancelled"
		case "confirmed":
			allowed = body.Status == "completed"
		case "cancelled", "completed":
			allowed = false
		default:
			allowed = false
		}
		if !allowed {
			c.JSON(400, gin.H{"error": "Cannot change status from " + currentStatus + " to " + body.Status})
			return
		}
		db.DB.Exec("UPDATE orders SET status = ?, updated_at = unixepoch() WHERE id = ?", body.Status, idStr)
		auditLog(uid, "order_status_changed", "order", idStr, currentStatus+" -> "+body.Status)
		BroadcastToUser(buyer, "order:status", gin.H{"order_id": idStr, "status": body.Status})
		BroadcastToUser(seller, "order:status", gin.H{"order_id": idStr, "status": body.Status})
	}
	if body.InstallmentPlan == "requested" || body.InstallmentPlan == "installments" {
		db.DB.Exec("UPDATE orders SET installment_plan = 'requested', updated_at = unixepoch() WHERE id = ?", idStr)
	}
	out := gin.H{"id": idStr}
	if body.Status != "" {
		out["status"] = body.Status
	}
	if body.InstallmentPlan != "" {
		out["installment_plan"] = "requested"
	}
	c.JSON(200, out)
}

func handleRemittanceCreate(c *gin.Context) {
	var body struct {
		ToIdentifier string  `json:"to_identifier"`
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
	}
	if c.ShouldBindJSON(&body) != nil {
		c.JSON(400, gin.H{"error": "to_identifier and amount required"})
		return
	}
	toID := strings.TrimSpace(body.ToIdentifier)
	if toID == "" {
		c.JSON(400, gin.H{"error": "to_identifier required"})
		return
	}
	if body.Amount <= 0 {
		c.JSON(400, gin.H{"error": "amount must be positive"})
		return
	}
	if body.Amount > 1e12 {
		c.JSON(400, gin.H{"error": "amount too large"})
		return
	}
	currency := strings.TrimSpace(body.Currency)
	if len(currency) > 10 {
		currency = currency[:10]
	}
	h, err := RemittanceCreate(getUserID(c), toID, body.Amount, currency)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create remittance request"})
		return
	}
	c.JSON(201, h)
}

func handleRemittancesMy(c *gin.Context) {
	list, err := RemittanceListMy(getUserID(c))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list remittances"})
		return
	}
	c.JSON(200, list)
}

func handleConversationsList(c *gin.Context) {
	c.JSON(200, ConversationsList(getUserID(c)))
}

func handleConversationsUnreadCount(c *gin.Context) {
	n := UnreadCount(getUserID(c))
	c.JSON(200, gin.H{"unread": n})
}

func handleConversationGet(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "unread-count" {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}
	cid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid conversation id"})
		return
	}
	h, err := ConversationGet(cid, getUserID(c))
	if err != nil {
		if errors.Is(err, ErrConvForbidden) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}
	c.JSON(200, h)
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
	if getUserID(c) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
	var otherID int64
	if db.DB.QueryRow("SELECT user_id FROM conversation_participants WHERE conversation_id = ? AND user_id != ?", cid, getUserID(c)).Scan(&otherID) == nil && otherID != 0 {
		BroadcastToUser(otherID, "relay:message", gin.H{"conversation_id": cid, "message_id": h["id"]})
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

// Vault API v1 (ARCHITECTURE-V4)
func handleVaultUploadURL(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Use POST /api/v1/vault/files with multipart for now; pre-signed URLs in Phase 3"})
}

func handleVaultCompleteUpload(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Pre-signed flow not implemented yet"})
}

func handleVaultDownloadURL(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Use GET /api/v1/vault/files/:id/download for now"})
}

func handleVaultListFiles(c *gin.Context) {
	uid := getUserID(c)
	folderID := c.Query("folder_id")
	var rows *sql.Rows
	var err error
	if folderID == "" {
		rows, err = db.DB.Query("SELECT id, name, size_bytes, mime_type, folder_id, created_at, updated_at FROM vault_files WHERE user_id = ? ORDER BY updated_at DESC", uid)
	} else {
		fid, _ := strconv.ParseInt(folderID, 10, 64)
		if fid == 0 {
			rows, err = db.DB.Query("SELECT id, name, size_bytes, mime_type, folder_id, created_at, updated_at FROM vault_files WHERE user_id = ? AND folder_id IS NULL ORDER BY updated_at DESC", uid)
		} else {
			rows, err = db.DB.Query("SELECT id, name, size_bytes, mime_type, folder_id, created_at, updated_at FROM vault_files WHERE user_id = ? AND folder_id = ? ORDER BY updated_at DESC", uid, fid)
		}
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list files"})
		return
	}
	defer rows.Close()
	var files []gin.H
	for rows.Next() {
		var id, sizeBytes int64
		var name, mimeType sql.NullString
		var folderIDNull sql.NullInt64
		var createdAt, updatedAt int64
		if err := rows.Scan(&id, &name, &sizeBytes, &mimeType, &folderIDNull, &createdAt, &updatedAt); err != nil {
			continue
		}
		f := gin.H{"id": id, "name": nullStrToString(name), "size_bytes": sizeBytes, "mime_type": nullStrToString(mimeType), "created_at": createdAt, "updated_at": updatedAt}
		if folderIDNull.Valid {
			f["folder_id"] = folderIDNull.Int64
		} else {
			f["folder_id"] = nil
		}
		files = append(files, f)
	}
	c.JSON(200, gin.H{"files": files})
}

func handleVaultListFolders(c *gin.Context) {
	uid := getUserID(c)
	parentID := c.Query("parent_id")
	var rows *sql.Rows
	var err error
	if parentID == "" {
		rows, err = db.DB.Query("SELECT id, name, parent_id, created_at, updated_at FROM vault_folders WHERE user_id = ? ORDER BY name", uid)
	} else {
		pid, _ := strconv.ParseInt(parentID, 10, 64)
		if pid == 0 {
			rows, err = db.DB.Query("SELECT id, name, parent_id, created_at, updated_at FROM vault_folders WHERE user_id = ? AND parent_id IS NULL ORDER BY name", uid)
		} else {
			rows, err = db.DB.Query("SELECT id, name, parent_id, created_at, updated_at FROM vault_folders WHERE user_id = ? AND parent_id = ? ORDER BY name", uid, pid)
		}
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list folders"})
		return
	}
	defer rows.Close()
	var folders []gin.H
	for rows.Next() {
		var id int64
		var name string
		var parentIDNull sql.NullInt64
		var createdAt, updatedAt int64
		if err := rows.Scan(&id, &name, &parentIDNull, &createdAt, &updatedAt); err != nil {
			continue
		}
		f := gin.H{"id": id, "name": name, "created_at": createdAt, "updated_at": updatedAt}
		if parentIDNull.Valid {
			f["parent_id"] = parentIDNull.Int64
		} else {
			f["parent_id"] = nil
		}
		folders = append(folders, f)
	}
	c.JSON(200, gin.H{"folders": folders})
}

func handleVaultCreateFolder(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		Name     string `json:"name"`
		ParentID *int64 `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" {
		c.JSON(400, gin.H{"error": "name required"})
		return
	}
	var parentID interface{}
	if body.ParentID != nil && *body.ParentID > 0 {
		var exists int
		if db.DB.QueryRow("SELECT 1 FROM vault_folders WHERE id = ? AND user_id = ?", *body.ParentID, uid).Scan(&exists) != nil {
			c.JSON(400, gin.H{"error": "parent folder not found"})
			return
		}
		parentID = *body.ParentID
	}
	res, err := db.DB.Exec("INSERT INTO vault_folders (user_id, name, parent_id, created_at, updated_at) VALUES (?, ?, ?, unixepoch(), unixepoch())", uid, body.Name, parentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create folder"})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(201, gin.H{"id": id, "name": body.Name, "parent_id": body.ParentID, "created_at": time.Now().Unix(), "updated_at": time.Now().Unix()})
}

func handleVaultDeleteFolder(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	res, err := db.DB.Exec("DELETE FROM vault_folders WHERE id = ? AND user_id = ?", id, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete folder"})
		return
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		c.JSON(404, gin.H{"error": "Folder not found"})
		return
	}
	// Orphan files in this folder (set folder_id = NULL)
	db.DB.Exec("UPDATE vault_files SET folder_id = NULL, updated_at = unixepoch() WHERE folder_id = ? AND user_id = ?", id, uid)
	c.JSON(200, gin.H{"ok": true})
}

func handleVaultUploadFile(c *gin.Context) {
	uid := getUserID(c)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "multipart form required"})
		return
	}
	files := form.File["file"]
	if len(files) == 0 {
		files = form.File["files"]
	}
	if len(files) == 0 {
		c.JSON(400, gin.H{"error": "file required"})
		return
	}
	fh := files[0]
	if fh.Size > cfg.MaxFileSize {
		c.JSON(400, gin.H{"error": "file too large"})
		return
	}
	folderIDRaw := form.Value["folder_id"]
	var folderID interface{}
	if len(folderIDRaw) > 0 && folderIDRaw[0] != "" {
		if fid, err := strconv.ParseInt(folderIDRaw[0], 10, 64); err == nil && fid > 0 {
			var exists int
			if db.DB.QueryRow("SELECT 1 FROM vault_folders WHERE id = ? AND user_id = ?", fid, uid).Scan(&exists) == nil {
				folderID = fid
			}
		}
	}
	// storage: uploads/vault/{user_id}/{file_id}
	res, err := db.DB.Exec("INSERT INTO vault_files (user_id, name, size_bytes, mime_type, storage_path, folder_id, created_at, updated_at) VALUES (?, ?, ?, ?, '', ?, unixepoch(), unixepoch())", uid, fh.Filename, fh.Size, fh.Header.Get("Content-Type"), folderID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create file record"})
		return
	}
	fileID, _ := res.LastInsertId()
	storagePath := filepath.Join(cfg.UploadDir, "vault", strconv.FormatInt(uid, 10), strconv.FormatInt(fileID, 10))
	if err := os.MkdirAll(filepath.Dir(storagePath), 0755); err != nil {
		db.DB.Exec("DELETE FROM vault_files WHERE id = ?", fileID)
		c.JSON(500, gin.H{"error": "Failed to create storage dir"})
		return
	}
	if err := c.SaveUploadedFile(fh, storagePath); err != nil {
		db.DB.Exec("DELETE FROM vault_files WHERE id = ?", fileID)
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}
	db.DB.Exec("UPDATE vault_files SET storage_path = ? WHERE id = ?", storagePath, fileID)
	vaultSearchIndexInsert(uid, fileID, fh.Filename)
	c.JSON(201, gin.H{"id": fileID, "name": fh.Filename, "size_bytes": fh.Size, "mime_type": fh.Header.Get("Content-Type"), "folder_id": folderID, "created_at": time.Now().Unix(), "updated_at": time.Now().Unix()})
}

func vaultSearchIndexInsert(userID, fileID int64, name string) {
	name = strings.ToLower(name)
	var buf []rune
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			buf = append(buf, r)
		} else if len(buf) >= 2 {
			word := string(buf)
			buf = nil
			h := sha256.Sum256([]byte(word))
			termHash := hex.EncodeToString(h[:])
			db.DB.Exec("INSERT INTO vault_search_index (user_id, file_id, term_hash, weight) VALUES (?, ?, ?, 1)", userID, fileID, termHash)
		} else {
			buf = nil
		}
	}
	if len(buf) >= 2 {
		h := sha256.Sum256([]byte(string(buf)))
		termHash := hex.EncodeToString(h[:])
		db.DB.Exec("INSERT INTO vault_search_index (user_id, file_id, term_hash, weight) VALUES (?, ?, ?, 1)", userID, fileID, termHash)
	}
}

func handleVaultGetFile(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var name, mimeType, storagePath string
	var sizeBytes, createdAt, updatedAt int64
	var folderIDNull sql.NullInt64
	err = db.DB.QueryRow("SELECT name, size_bytes, mime_type, storage_path, folder_id, created_at, updated_at FROM vault_files WHERE id = ? AND user_id = ?", id, uid).Scan(&name, &sizeBytes, &mimeType, &storagePath, &folderIDNull, &createdAt, &updatedAt)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	out := gin.H{"id": id, "name": name, "size_bytes": sizeBytes, "mime_type": mimeType, "created_at": createdAt, "updated_at": updatedAt}
	if folderIDNull.Valid {
		out["folder_id"] = folderIDNull.Int64
	} else {
		out["folder_id"] = nil
	}
	c.JSON(200, out)
}

func handleVaultDownloadFile(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var name, storagePath string
	err = db.DB.QueryRow("SELECT name, storage_path FROM vault_files WHERE id = ? AND user_id = ?", id, uid).Scan(&name, &storagePath)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	if _, err := os.Stat(storagePath); err != nil {
		c.JSON(404, gin.H{"error": "File not found on disk"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=\""+name+"\"")
	c.File(storagePath)
}

func handleVaultDeleteFile(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var storagePath string
	if db.DB.QueryRow("SELECT storage_path FROM vault_files WHERE id = ? AND user_id = ?", id, uid).Scan(&storagePath) != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	db.DB.Exec("DELETE FROM vault_files WHERE id = ?", id)
	os.Remove(storagePath)
	c.JSON(200, gin.H{"ok": true})
}

func handleAuthSessionsList(c *gin.Context) {
	uid := getUserID(c)
	rows, err := db.DB.Query("SELECT id, device_name, created_at, expires_at FROM sessions WHERE user_id = ? AND expires_at > ? ORDER BY created_at DESC", uid, time.Now().Unix())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to list sessions"})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id int64
		var deviceName string
		var createdAt, expiresAt int64
		if rows.Scan(&id, &deviceName, &createdAt, &expiresAt) != nil {
			continue
		}
		list = append(list, gin.H{"id": id, "device_name": deviceName, "created_at": createdAt, "expires_at": expiresAt})
	}
	c.JSON(200, gin.H{"sessions": list})
}

func handleAuthSessionDelete(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	res, err := db.DB.Exec("DELETE FROM sessions WHERE id = ? AND user_id = ?", id, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to delete"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(404, gin.H{"error": "session not found"})
		return
	}
	auditLog(uid, "session.revoked", "session", strconv.FormatInt(id, 10), "")
	c.JSON(200, gin.H{"ok": true})
}

func handleAuthDevicesList(c *gin.Context) {
	uid := getUserID(c)
	rows, err := db.DB.Query("SELECT id, name, last_used, created_at FROM devices WHERE user_id = ? ORDER BY created_at DESC", uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to list devices"})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id int64
		var name string
		var lastUsed sql.NullInt64
		var createdAt int64
		if rows.Scan(&id, &name, &lastUsed, &createdAt) != nil {
			continue
		}
		lu := interface{}(nil)
		if lastUsed.Valid {
			lu = lastUsed.Int64
		}
		list = append(list, gin.H{"id": id, "name": name, "last_used": lu, "created_at": createdAt})
	}
	c.JSON(200, gin.H{"devices": list})
}

func handleAuthDeviceDelete(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	res, err := db.DB.Exec("DELETE FROM devices WHERE id = ? AND user_id = ?", id, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to delete"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(404, gin.H{"error": "device not found"})
		return
	}
	auditLog(uid, "device.removed", "device", strconv.FormatInt(id, 10), "")
	c.JSON(200, gin.H{"ok": true})
}

func handleRecoveryGenerate(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		RecoveryHash string `json:"recoveryHash"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.RecoveryHash == "" {
		c.JSON(400, gin.H{"error": "recoveryHash required"})
		return
	}
	_, err := db.DB.Exec(
		"INSERT OR REPLACE INTO user_recovery (user_id, recovery_hash, created_at) VALUES (?, ?, ?)",
		uid, body.RecoveryHash, time.Now().Unix(),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save recovery"})
		return
	}
	auditLog(uid, "recovery.generated", "user_recovery", "", "")
	c.JSON(200, gin.H{"ok": true})
}

func handleRecoveryVerify(c *gin.Context) {
	var body struct {
		RecoveryHash string `json:"recoveryHash"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.RecoveryHash == "" {
		c.JSON(400, gin.H{"error": "recoveryHash required"})
		return
	}
	var userID int64
	if db.DB.QueryRow("SELECT user_id FROM user_recovery WHERE recovery_hash = ?", body.RecoveryHash).Scan(&userID) != nil {
		c.JSON(400, gin.H{"error": "invalid recovery phrase"})
		return
	}
	c.JSON(200, gin.H{"valid": true, "userId": userID})
}

func handleRecoveryRestore(c *gin.Context) {
	var body struct {
		RecoveryHash string `json:"recoveryHash"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.RecoveryHash == "" {
		c.JSON(400, gin.H{"error": "recoveryHash required"})
		return
	}
	var userID int64
	if db.DB.QueryRow("SELECT user_id FROM user_recovery WHERE recovery_hash = ?", body.RecoveryHash).Scan(&userID) != nil {
		c.JSON(400, gin.H{"error": "invalid recovery phrase"})
		return
	}
	db.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	sessionID, exp := createSession(userID, "recovery")
	token, _ := pqc.SignTokenWithSession(cfg.PQCPrivateKey, userID, sessionID, exp)
	auditLog(userID, "recovery.restore", "user_recovery", "", "")
	c.JSON(200, gin.H{"token": token, "user_id": userID})
}

func nullStrToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
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
