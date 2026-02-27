// Auth service layer: register and login. Handlers stay thin and map service errors to HTTP.
package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"omnixius-api/db"
	"omnixius-api/pqc"

	"github.com/gin-gonic/gin"
)

// createSession inserts a session row and returns (sessionID, expiresAt). Caller uses SignTokenWithSession.
func createSession(userID int64, deviceName string) (sessionID int64, expiresAt time.Time) {
	exp := time.Now().Add(7 * 24 * time.Hour)
	res, err := db.DB.Exec(
		"INSERT INTO sessions (user_id, device_name, expires_at) VALUES (?, ?, ?)",
		userID, deviceName, exp.Unix(),
	)
	if err != nil {
		return 0, exp
	}
	sid, _ := res.LastInsertId()
	return sid, exp
}

var (
	ErrEmailExists         = errors.New("email already registered")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrRegistrationFailed  = errors.New("registration failed")
)

// AuthRegister creates a user and returns user map and token. Caller must validate input length and min password.
func AuthRegister(email, password, name string) (user gin.H, token string, err error) {
	email = strings.TrimSpace(strings.ToLower(email))
	var id int64
	if db.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&id) == nil {
		return nil, "", ErrEmailExists
	}
	hash := hashPasswordArgon2(password)
	verifyToken := make([]byte, 32)
	rand.Read(verifyToken)
	// email_verified = 1 until we have mail sending; then set 0 and send confirm link
	res, err := db.DB.Exec(
		"INSERT INTO users (email, password_hash, name, role, email_verify_token, email_verified) VALUES (?, ?, ?, 'user', ?, 1)",
		email, hash, nullStr(name), hex.EncodeToString(verifyToken),
	)
	if err != nil {
		return nil, "", ErrRegistrationFailed
	}
	id, _ = res.LastInsertId()
	sessionID, exp := createSession(id, "web")
	token, _ = pqc.SignTokenWithSession(cfg.PQCPrivateKey, id, sessionID, exp)
	user = gin.H{"id": id, "email": email, "role": "user", "name": name}
	return user, token, nil
}

// AuthLogin returns user map and token. Caller must enforce rate limit.
func AuthLogin(email, password string) (user gin.H, token string, err error) {
	email = strings.TrimSpace(strings.ToLower(email))
	var id int64
	var hash string
	var role, name string
	var avatar sql.NullString
	var emailVerified, phoneVerified int
	err = db.DB.QueryRow(
		"SELECT id, password_hash, role, name, avatar_path, COALESCE(email_verified, 0), COALESCE(phone_verified, 0) FROM users WHERE email = ?",
		email,
	).Scan(&id, &hash, &role, &name, &avatar, &emailVerified, &phoneVerified)
	if err != nil || !checkPassword(hash, password) {
		return nil, "", ErrInvalidCredentials
	}
	verified := emailVerified == 1 || phoneVerified == 1
	user = gin.H{
		"id": id, "email": email, "role": role, "name": name,
		"avatar_path": avatar.String, "email_verified": emailVerified == 1, "phone_verified": phoneVerified == 1, "verified": verified,
	}
	sessionID, exp := createSession(id, "web")
	token, _ = pqc.SignTokenWithSession(cfg.PQCPrivateKey, id, sessionID, exp)
	return user, token, nil
}
