// WebAuthn (Passkeys) handlers for ARCHITECTURE-V4 Phase 1.
package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"omnixius-api/db"
	"omnixius-api/pqc"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

var webauthnInstance *webauthn.WebAuthn

func initWebAuthn() error {
	config := &webauthn.Config{
		RPID:          cfg.WebAuthnRPID,
		RPDisplayName: cfg.WebAuthnRPDisplayName,
		RPOrigins:     cfg.WebAuthnRPOrigins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementPreferred,
			UserVerification: protocol.VerificationPreferred,
		},
		AttestationPreference: protocol.PreferNoAttestation,
	}
	var err error
	webauthnInstance, err = webauthn.New(config)
	return err
}

// webauthnUser implements webauthn.User for a DB user.
type webauthnUser struct {
	id          int64
	email       string
	displayName string
	credentials []webauthn.Credential
}

func (u webauthnUser) WebAuthnID() []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(u.id >> (i * 8))
	}
	return b
}

func (u webauthnUser) WebAuthnName() string       { return u.email }
func (u webauthnUser) WebAuthnDisplayName() string { return u.displayName }
func (u webauthnUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

func loadWebAuthnCredentials(userID int64) ([]webauthn.Credential, error) {
	rows, err := db.DB.Query(
		"SELECT credential_json FROM webauthn_credentials WHERE user_id = ? ORDER BY id",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var creds []webauthn.Credential
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var c webauthn.Credential
		if json.Unmarshal([]byte(raw), &c) != nil {
			continue
		}
		creds = append(creds, c)
	}
	return creds, nil
}

func saveWebAuthnCredential(userID int64, cred *webauthn.Credential) error {
	raw, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	idB64 := base64.RawURLEncoding.EncodeToString(cred.ID)
	_, err = db.DB.Exec(
		"INSERT INTO webauthn_credentials (user_id, credential_id_base64, credential_json, updated_at) VALUES (?, ?, ?, ?)",
		userID, idB64, string(raw), time.Now().Unix(),
	)
	return err
}

func webauthnSaveSession(session *webauthn.SessionData) (string, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()
	_, err = db.DB.Exec(
		"INSERT INTO webauthn_sessions (id, session_data, created_at) VALUES (?, ?, ?)",
		id, string(data), time.Now().Unix(),
	)
	return id, err
}

func webauthnLoadSession(sessionID string) (*webauthn.SessionData, error) {
	var raw string
	err := db.DB.QueryRow("SELECT session_data FROM webauthn_sessions WHERE id = ?", sessionID).Scan(&raw)
	if err != nil {
		return nil, err
	}
	var s webauthn.SessionData
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func webauthnDeleteSession(sessionID string) {
	db.DB.Exec("DELETE FROM webauthn_sessions WHERE id = ?", sessionID)
}

func handlePasskeyRegisterBegin(c *gin.Context) {
	if webauthnInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "WebAuthn not configured"})
		return
	}
	var body struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email required"})
		return
	}
	email := strings.TrimSpace(strings.ToLower(body.Email))
	if email == "" || !isValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid email required"})
		return
	}
	var existingID int64
	if db.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&existingID) == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered; use login or add passkey to existing account"})
		return
	}
	// Create user with placeholder password (passkey-only account)
	hash := hashPasswordArgon2(string(uuid.New().String())) // unguessable
	res, err := db.DB.Exec(
		"INSERT INTO users (email, password_hash, name, role) VALUES (?, ?, ?, 'user')",
		email, hash, nullStr(body.Name),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}
	userID, _ := res.LastInsertId()
	creds, _ := loadWebAuthnCredentials(userID)
	u := webauthnUser{id: userID, email: email, displayName: body.Name, credentials: creds}
	creation, session, err := webauthnInstance.BeginRegistration(u)
	if err != nil {
		db.DB.Exec("DELETE FROM users WHERE id = ?", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start registration"})
		return
	}
	sessionID, err := webauthnSaveSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session save failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"options":    creation,
	})
}

func handlePasskeyRegisterComplete(c *gin.Context) {
	if webauthnInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "WebAuthn not configured"})
		return
	}
	sessionID := c.GetHeader("X-WebAuthn-Session")
	if sessionID == "" {
		sessionID = c.Query("session_id")
	}
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-WebAuthn-Session or session_id required"})
		return
	}
	session, err := webauthnLoadSession(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired session"})
		return
	}
	defer webauthnDeleteSession(sessionID)
	var userID int64
	for i := 0; i < 8 && i < len(session.UserID); i++ {
		userID |= int64(session.UserID[i]) << (i * 8)
	}
	var email, name string
	if db.DB.QueryRow("SELECT email, name FROM users WHERE id = ?", userID).Scan(&email, &name) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	creds, _ := loadWebAuthnCredentials(userID)
	u := webauthnUser{id: userID, email: email, displayName: name, credentials: creds}
	cred, err := webauthnInstance.FinishRegistration(u, *session, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "registration failed: " + err.Error()})
		return
	}
	if err := saveWebAuthnCredential(userID, cred); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save credential"})
		return
	}
	authSid, authExp := createSession(userID, "passkey")
	token, _ := pqc.SignTokenWithSession(cfg.PQCPrivateKey, userID, authSid, authExp)
	c.JSON(http.StatusOK, gin.H{
		"user":  gin.H{"id": userID, "email": email, "role": "user", "name": name},
		"token": token,
	})
}

func handlePasskeyLoginBegin(c *gin.Context) {
	if webauthnInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "WebAuthn not configured"})
		return
	}
	var body struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email required"})
		return
	}
	email := strings.TrimSpace(strings.ToLower(body.Email))
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email required"})
		return
	}
	var userID int64
	var name sql.NullString
	if db.DB.QueryRow("SELECT id, name FROM users WHERE email = ?", email).Scan(&userID, &name) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no account with this email"})
		return
	}
	creds, err := loadWebAuthnCredentials(userID)
	if err != nil || len(creds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no passkey registered for this account; use password login"})
		return
	}
	displayName := name.String
	if displayName == "" {
		displayName = email
	}
	u := webauthnUser{id: userID, email: email, displayName: displayName, credentials: creds}
	assertion, session, err := webauthnInstance.BeginLogin(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start login"})
		return
	}
	sessionID, err := webauthnSaveSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session save failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"options":    assertion,
	})
}

func handlePasskeyLoginComplete(c *gin.Context) {
	if webauthnInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "WebAuthn not configured"})
		return
	}
	sessionID := c.GetHeader("X-WebAuthn-Session")
	if sessionID == "" {
		sessionID = c.Query("session_id")
	}
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-WebAuthn-Session or session_id required"})
		return
	}
	session, err := webauthnLoadSession(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired session"})
		return
	}
	defer webauthnDeleteSession(sessionID)
	var userID int64
	for i := 0; i < 8 && i < len(session.UserID); i++ {
		userID |= int64(session.UserID[i]) << (i * 8)
	}
	var email, role string
	var name, avatar sql.NullString
	if db.DB.QueryRow("SELECT email, role, name, avatar_path FROM users WHERE id = ?", userID).Scan(&email, &role, &name, &avatar) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	creds, _ := loadWebAuthnCredentials(userID)
	u := webauthnUser{id: userID, email: email, displayName: name.String, credentials: creds}
	_, err = webauthnInstance.FinishLogin(u, *session, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed: " + err.Error()})
		return
	}
	authSid, authExp := createSession(userID, "passkey")
	token, _ := pqc.SignTokenWithSession(cfg.PQCPrivateKey, userID, authSid, authExp)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id": userID, "email": email, "role": role, "name": name.String,
			"avatar_path": avatar.String,
		},
		"token": token,
	})
}
