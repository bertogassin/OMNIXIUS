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
	res, err := db.DB.Exec(
		"INSERT INTO users (email, password_hash, name, role, email_verify_token) VALUES (?, ?, ?, 'user', ?)",
		email, hash, nullStr(name), hex.EncodeToString(verifyToken),
	)
	if err != nil {
		return nil, "", ErrRegistrationFailed
	}
	id, _ = res.LastInsertId()
	token, _ = pqc.SignToken(cfg.PQCPrivateKey, id, time.Now().Add(7*24*time.Hour))
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
	var verified int
	err = db.DB.QueryRow(
		"SELECT id, password_hash, role, name, avatar_path, email_verified FROM users WHERE email = ?",
		email,
	).Scan(&id, &hash, &role, &name, &avatar, &verified)
	if err != nil || !checkPassword(hash, password) {
		return nil, "", ErrInvalidCredentials
	}
	user = gin.H{
		"id": id, "email": email, "role": role, "name": name,
		"avatar_path": avatar.String, "email_verified": verified == 1,
	}
	token, _ = pqc.SignToken(cfg.PQCPrivateKey, id, time.Now().Add(7*24*time.Hour))
	return user, token, nil
}
