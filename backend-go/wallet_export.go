// Wallet backup: encrypted export/import (WHAT-WE-TAKE). AES-256-GCM, key from password (Argon2).
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

const walletExportVersion = 1

type walletExportPayload struct {
	Version   int       `json:"version"`
	UserID    int64     `json:"user_id"`
	ExportedAt int64    `json:"exported_at"`
	Balances  []gin.H   `json:"balances"`
	Addresses []gin.H   `json:"addresses"`
}

func handleWalletExport(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
		return
	}
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate salt"})
		return
	}
	key := argon2.IDKey([]byte(body.Password), salt, 1, 64*1024, 4, 32)

	var balances []gin.H
	rows, _ := db.DB.Query("SELECT currency, amount, hold_amount, updated_at FROM wallet_balances WHERE user_id = ?", uid)
	if rows != nil {
		for rows.Next() {
			var currency string
			var amount, holdAmount int64
			var updatedAt sql.NullInt64
			if rows.Scan(&currency, &amount, &holdAmount, &updatedAt) == nil {
				balances = append(balances, gin.H{"currency": currency, "amount": amount, "hold_amount": holdAmount, "updated_at": updatedAt.Int64})
			}
		}
		rows.Close()
	}
	var addresses []gin.H
	rows2, _ := db.DB.Query("SELECT id, currency, address, network, created_at, last_used_at FROM wallet_deposit_addresses WHERE user_id = ?", uid)
	if rows2 != nil {
		for rows2.Next() {
			var id int64
			var currency, address, network string
			var createdAt, lastUsedAt sql.NullInt64
			if rows2.Scan(&id, &currency, &address, &network, &createdAt, &lastUsedAt) == nil {
				addresses = append(addresses, gin.H{"id": id, "currency": currency, "address": address, "network": network, "created_at": createdAt.Int64, "last_used_at": lastUsedAt.Int64})
			}
		}
		rows2.Close()
	}
	payload := walletExportPayload{
		Version:    walletExportVersion,
		UserID:     uid,
		ExportedAt: time.Now().Unix(),
		Balances:   balances,
		Addresses:  addresses,
	}
	plain, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export failed"})
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export failed"})
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export failed"})
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export failed"})
		return
	}
	ciphertext := gcm.Seal(nonce, nonce, plain, nil)
	c.JSON(http.StatusOK, gin.H{
		"export":  base64.StdEncoding.EncodeToString(ciphertext),
		"salt":    base64.StdEncoding.EncodeToString(salt),
		"version": walletExportVersion,
	})
}

func handleWalletImport(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		Export   string `json:"export"`
		Salt     string `json:"salt"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Export == "" || body.Salt == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "export, salt and password required"})
		return
	}
	salt, err := base64.StdEncoding.DecodeString(body.Salt)
	if err != nil || len(salt) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid salt"})
		return
	}
	ciphertext, err := base64.StdEncoding.DecodeString(body.Export)
	if err != nil || len(ciphertext) < 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid export"})
		return
	}
	key := argon2.IDKey([]byte(body.Password), salt, 1, 64*1024, 4, 32)
	block, err := aes.NewCipher(key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "decrypt failed"})
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "decrypt failed"})
		return
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid export"})
		return
	}
	plain, err := gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong password or corrupted export"})
		return
	}
	var payload walletExportPayload
	if err := json.Unmarshal(plain, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid export data"})
		return
	}
	if payload.UserID != uid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "export belongs to another user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"verified": true,
		"version":  payload.Version,
		"exported_at": payload.ExportedAt,
		"balances_count":  len(payload.Balances),
		"addresses_count": len(payload.Addresses),
	})
}
