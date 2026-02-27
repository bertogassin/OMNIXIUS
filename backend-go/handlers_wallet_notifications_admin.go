// Wallet (§15), Notifications (§16), Vault search (§17), Admin (§18)
package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

// --- Wallet ---
func handleWalletBalances(c *gin.Context) {
	uid := getUserID(c)
	rows, err := db.DB.Query(
		"SELECT currency, amount, hold_amount, updated_at FROM wallet_balances WHERE user_id = ?",
		uid,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load balances"})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var currency string
		var amount, holdAmount int64
		var updatedAt sql.NullInt64
		if rows.Scan(&currency, &amount, &holdAmount, &updatedAt) != nil {
			continue
		}
		list = append(list, gin.H{
			"currency":    currency,
			"amount":      amount,
			"hold_amount": holdAmount,
			"available":   amount - holdAmount,
			"updated_at":  updatedAt.Int64,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"balances": list})
}

func handleWalletTransactions(c *gin.Context) {
	uid := getUserID(c)
	limit := 50
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	rows, err := db.DB.Query(
		"SELECT id, type, currency, amount, fee, status, reference_id, created_at FROM wallet_transactions WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		uid, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load transactions"})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id int64
		var txType, currency, status, refID string
		var amount, fee int64
		var createdAt sql.NullInt64
		if rows.Scan(&id, &txType, &currency, &amount, &fee, &status, &refID, &createdAt) != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "type": txType, "currency": currency,
			"amount": amount, "fee": fee, "status": status,
			"reference_id": refID, "created_at": createdAt.Int64,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"transactions": list})
}

func handleWalletTransfer(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		ToUserID int64  `json:"to_user_id"`
		Currency string `json:"currency"`
		Amount   int64  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ToUserID <= 0 || body.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to_user_id, currency, amount (positive) required"})
		return
	}
	if body.Currency == "" {
		body.Currency = "USD"
	}
	if body.ToUserID == uid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot transfer to self"})
		return
	}
	// Ensure sender has wallet_balances row and enough balance
	var amount, holdAmount int64
	err := db.DB.QueryRow(
		"SELECT amount, hold_amount FROM wallet_balances WHERE user_id = ? AND currency = ?",
		uid, body.Currency,
	).Scan(&amount, &holdAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}
	if amount-holdAmount < body.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}
	now := time.Now().Unix()
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transfer failed"})
		return
	}
	defer tx.Rollback()
	_, err = tx.Exec(
		"UPDATE wallet_balances SET amount = amount - ?, updated_at = ? WHERE user_id = ? AND currency = ?",
		body.Amount, now, uid, body.Currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transfer failed"})
		return
	}
	_, err = tx.Exec(
		"INSERT INTO wallet_balances (user_id, currency, amount, hold_amount, updated_at) VALUES (?, ?, ?, 0, ?) ON CONFLICT(user_id, currency) DO UPDATE SET amount = amount + ?, updated_at = ?",
		body.ToUserID, body.Currency, body.Amount, now, body.Amount, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transfer failed"})
		return
	}
	_, _ = tx.Exec(
		"INSERT INTO wallet_transactions (user_id, type, currency, amount, fee, status, reference_id, created_at, completed_at) VALUES (?, 'transfer_out', ?, ?, 0, 'completed', ?, ?, ?)",
		uid, body.Currency, -body.Amount, strconv.FormatInt(body.ToUserID, 10), now, now,
	)
	_, _ = tx.Exec(
		"INSERT INTO wallet_transactions (user_id, type, currency, amount, fee, status, reference_id, created_at, completed_at) VALUES (?, 'transfer_in', ?, ?, 0, 'completed', ?, ?, ?)",
		body.ToUserID, body.Currency, body.Amount, strconv.FormatInt(uid, 10), now, now,
	)
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transfer failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleWalletBalanceByCurrency(c *gin.Context) {
	uid := getUserID(c)
	currency := c.Param("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency required"})
		return
	}
	var amount, holdAmount int64
	var updatedAt sql.NullInt64
	err := db.DB.QueryRow(
		"SELECT amount, hold_amount, updated_at FROM wallet_balances WHERE user_id = ? AND currency = ?",
		uid, currency,
	).Scan(&amount, &holdAmount, &updatedAt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"currency": currency, "amount": 0, "hold_amount": 0, "available": 0, "updated_at": 0,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"currency": currency, "amount": amount, "hold_amount": holdAmount,
		"available": amount - holdAmount, "updated_at": updatedAt.Int64,
	})
}

func handleWalletTransactionByID(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var txType, currency, status, refID string
	var amount, fee int64
	var createdAt, completedAt sql.NullInt64
	err = db.DB.QueryRow(
		"SELECT type, currency, amount, fee, status, reference_id, created_at, completed_at FROM wallet_transactions WHERE id = ? AND user_id = ?",
		id, uid,
	).Scan(&txType, &currency, &amount, &fee, &status, &refID, &createdAt, &completedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id, "type": txType, "currency": currency, "amount": amount, "fee": fee,
		"status": status, "reference_id": refID, "created_at": createdAt.Int64, "completed_at": completedAt.Int64,
	})
}

func handleWalletTransferVerify(c *gin.Context) {
	var body struct {
		ToUserID int64 `json:"to_user_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ToUserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to_user_id required"})
		return
	}
	var exists int
	err := db.DB.QueryRow("SELECT 1 FROM users WHERE id = ?", body.ToUserID).Scan(&exists)
	if err != nil || exists != 1 {
		c.JSON(http.StatusOK, gin.H{"valid": false, "error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"valid": true})
}

func handleWalletDepositAddressesList(c *gin.Context) {
	uid := getUserID(c)
	rows, err := db.DB.Query(
		"SELECT id, currency, address, network, created_at, last_used_at FROM wallet_deposit_addresses WHERE user_id = ? ORDER BY created_at DESC",
		uid,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"addresses": []gin.H{}})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id int64
		var currency, address, network string
		var createdAt sql.NullInt64
		var lastUsedAt sql.NullInt64
		if rows.Scan(&id, &currency, &address, &network, &createdAt, &lastUsedAt) != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "currency": currency, "address": address, "network": network,
			"created_at": createdAt.Int64, "last_used_at": lastUsedAt.Int64,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": list})
}

func handleWalletDepositAddressCreate(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		Currency string `json:"currency"`
		Network  string `json:"network"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency required"})
		return
	}
	if body.Network == "" {
		body.Network = "mainnet"
	}
	// Generate a placeholder address (real impl would call blockchain API)
	address := "0x" + strconv.FormatInt(uid, 10) + "_" + body.Currency + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	res, err := db.DB.Exec(
		"INSERT INTO wallet_deposit_addresses (user_id, currency, address, network) VALUES (?, ?, ?, ?)",
		uid, body.Currency, address, body.Network,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "address already exists or failed"})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id, "currency": body.Currency, "address": address, "network": body.Network})
}

func handleWalletHold(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		OrderID  *int64 `json:"order_id"`
		Currency string `json:"currency"`
		Amount   int64  `json:"amount"`
		ExpiresIn int   `json:"expires_in"` // seconds from now
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Currency == "" || body.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency and amount (positive) required"})
		return
	}
	if body.ExpiresIn <= 0 {
		body.ExpiresIn = 86400 * 7 // 7 days
	}
	expiresAt := time.Now().Unix() + int64(body.ExpiresIn)
	var amount, holdAmount int64
	err := db.DB.QueryRow(
		"SELECT amount, hold_amount FROM wallet_balances WHERE user_id = ? AND currency = ?",
		uid, body.Currency,
	).Scan(&amount, &holdAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}
	if amount-holdAmount < body.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}
	now := time.Now().Unix()
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hold failed"})
		return
	}
	defer tx.Rollback()
	var orderID interface{}
	if body.OrderID != nil {
		orderID = *body.OrderID
	}
	res, err := tx.Exec(
		"INSERT INTO wallet_holds (user_id, order_id, currency, amount, expires_at) VALUES (?, ?, ?, ?, ?)",
		uid, orderID, body.Currency, body.Amount, expiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hold failed"})
		return
	}
	holdID, _ := res.LastInsertId()
	_, err = tx.Exec(
		"UPDATE wallet_balances SET hold_amount = hold_amount + ?, updated_at = ? WHERE user_id = ? AND currency = ?",
		body.Amount, now, uid, body.Currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hold failed"})
		return
	}
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hold failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": holdID, "expires_at": expiresAt})
}

func handleWalletHoldRelease(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var holdUserID int64
	var currency string
	var amount int64
	var releasedAt sql.NullInt64
	err = db.DB.QueryRow(
		"SELECT user_id, currency, amount, released_at FROM wallet_holds WHERE id = ?",
		id,
	).Scan(&holdUserID, &currency, &amount, &releasedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if holdUserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if releasedAt.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hold already released"})
		return
	}
	now := time.Now().Unix()
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "release failed"})
		return
	}
	defer tx.Rollback()
	_, err = tx.Exec("UPDATE wallet_holds SET released_at = ? WHERE id = ?", now, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "release failed"})
		return
	}
	_, err = tx.Exec(
		"UPDATE wallet_balances SET hold_amount = hold_amount - ?, updated_at = ? WHERE user_id = ? AND currency = ?",
		amount, now, uid, currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "release failed"})
		return
	}
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "release failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleWalletHoldCapture(c *gin.Context) {
	uid := getUserID(c) // may be buyer or admin; we allow capture if caller is involved in order or admin
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		ToUserID int64 `json:"to_user_id"` // seller_id for trade
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ToUserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to_user_id required"})
		return
	}
	var holdUserID int64
	var currency string
	var amount int64
	var releasedAt sql.NullInt64
	err = db.DB.QueryRow(
		"SELECT user_id, currency, amount, released_at FROM wallet_holds WHERE id = ?",
		id,
	).Scan(&holdUserID, &currency, &amount, &releasedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if releasedAt.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hold already released or captured"})
		return
	}
	// Caller must be the user who owns the hold (buyer confirming delivery)
	if holdUserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	now := time.Now().Unix()
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "capture failed"})
		return
	}
	defer tx.Rollback()
	_, err = tx.Exec("UPDATE wallet_holds SET released_at = ? WHERE id = ?", now, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "capture failed"})
		return
	}
	_, err = tx.Exec(
		"UPDATE wallet_balances SET hold_amount = hold_amount - ?, updated_at = ? WHERE user_id = ? AND currency = ?",
		amount, now, holdUserID, currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "capture failed"})
		return
	}
	_, err = tx.Exec(
		"INSERT INTO wallet_balances (user_id, currency, amount, hold_amount, updated_at) VALUES (?, ?, ?, 0, ?) ON CONFLICT(user_id, currency) DO UPDATE SET amount = amount + ?, updated_at = ?",
		body.ToUserID, currency, amount, now, amount, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "capture failed"})
		return
	}
	_, _ = tx.Exec(
		"INSERT INTO wallet_transactions (user_id, type, currency, amount, fee, status, reference_id, created_at, completed_at) VALUES (?, 'payment', ?, ?, 0, 'completed', ?, ?, ?)",
		body.ToUserID, currency, amount, "hold:"+strconv.FormatInt(id, 10), now, now,
	)
	_, _ = tx.Exec(
		"INSERT INTO wallet_transactions (user_id, type, currency, amount, fee, status, reference_id, created_at, completed_at) VALUES (?, 'payment', ?, ?, 0, 'completed', ?, ?, ?)",
		holdUserID, currency, -amount, "hold:"+strconv.FormatInt(id, 10), now, now,
	)
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "capture failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// --- Notifications ---
func handleNotificationsSettingsGet(c *gin.Context) {
	uid := getUserID(c)
	var emailEnabled, pushEnabled int
	err := db.DB.QueryRow(
		"SELECT COALESCE(email_enabled, 1), COALESCE(push_enabled, 1) FROM notifications_user_settings WHERE user_id = ?",
		uid,
	).Scan(&emailEnabled, &pushEnabled)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"email_enabled": true, "push_enabled": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"email_enabled": emailEnabled == 1,
		"push_enabled":  pushEnabled == 1,
	})
}

func handleNotificationsSettingsPatch(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		EmailEnabled *bool `json:"email_enabled"`
		PushEnabled  *bool `json:"push_enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	now := time.Now().Unix()
	var emailEnabled, pushEnabled int
	_ = db.DB.QueryRow("SELECT COALESCE(email_enabled, 1), COALESCE(push_enabled, 1) FROM notifications_user_settings WHERE user_id = ?", uid).Scan(&emailEnabled, &pushEnabled)
	if body.EmailEnabled != nil {
		if *body.EmailEnabled {
			emailEnabled = 1
		} else {
			emailEnabled = 0
		}
	}
	if body.PushEnabled != nil {
		if *body.PushEnabled {
			pushEnabled = 1
		} else {
			pushEnabled = 0
		}
	}
	_, err := db.DB.Exec(
		"INSERT INTO notifications_user_settings (user_id, email_enabled, push_enabled, updated_at) VALUES (?, ?, ?, ?) ON CONFLICT(user_id) DO UPDATE SET email_enabled = excluded.email_enabled, push_enabled = excluded.push_enabled, updated_at = excluded.updated_at",
		uid, emailEnabled, pushEnabled, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleNotificationsHistory(c *gin.Context) {
	uid := getUserID(c)
	limit := 30
	rows, err := db.DB.Query(
		"SELECT id, type, channel, title, body, data, status, sent_at, created_at FROM notifications_queue WHERE user_id = ? ORDER BY created_at DESC LIMIT ?",
		uid, limit,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"notifications": []gin.H{}})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id int64
		var ntype, channel, title, body, status string
		var data sql.NullString
		var sentAt, createdAt sql.NullInt64
		if rows.Scan(&id, &ntype, &channel, &title, &body, &data, &status, &sentAt, &createdAt) != nil {
			continue
		}
		item := gin.H{
			"id": id, "type": ntype, "channel": channel,
			"title": title, "body": body, "status": status,
			"sent_at": sentAt.Int64, "created_at": createdAt.Int64,
		}
		if data.Valid && data.String != "" {
			item["data"] = data.String
		}
		list = append(list, item)
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"notifications": list})
}

func handleNotificationsHistoryByID(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var ntype, channel, title, body, status string
	var sentAt, createdAt sql.NullInt64
	err = db.DB.QueryRow(
		"SELECT type, channel, title, body, status, sent_at, created_at FROM notifications_queue WHERE id = ? AND user_id = ?",
		id, uid,
	).Scan(&ntype, &channel, &title, &body, &status, &sentAt, &createdAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	out := gin.H{"id": id, "type": ntype, "channel": channel, "title": title, "body": body, "status": status, "sent_at": sentAt.Int64, "created_at": createdAt.Int64}
	var readAt sql.NullInt64
	_ = db.DB.QueryRow("SELECT read_at FROM notifications_queue WHERE id = ? AND user_id = ?", id, uid).Scan(&readAt)
	if readAt.Valid {
		out["read_at"] = readAt.Int64
	} else {
		out["read_at"] = nil
	}
	c.JSON(http.StatusOK, out)
}

func handleNotificationsHistoryRead(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	now := time.Now().Unix()
	res, err := db.DB.Exec("UPDATE notifications_queue SET read_at = ? WHERE id = ? AND user_id = ?", now, id, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleNotificationsPushTokenCreate(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
		DeviceID *int64 `json:"device_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
		return
	}
	if body.Platform == "" {
		body.Platform = "web"
	}
	var deviceID interface{}
	if body.DeviceID != nil {
		deviceID = *body.DeviceID
	}
	res, err := db.DB.Exec(
		"INSERT INTO notifications_push_tokens (user_id, device_id, token, platform) VALUES (?, ?, ?, ?)",
		uid, deviceID, body.Token, body.Platform,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register token"})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id, "ok": true})
}

func handleNotificationsPushTokenDelete(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	res, err := db.DB.Exec("DELETE FROM notifications_push_tokens WHERE id = ? AND user_id = ?", id, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleNotificationsTest(c *gin.Context) {
	uid := getUserID(c)
	title := "Test notification"
	body := "This is a test notification from OMNIXIUS."
	now := time.Now().Unix()
	_, err := db.DB.Exec(
		"INSERT INTO notifications_queue (user_id, type, channel, title, body, status, scheduled_for) VALUES (?, ?, ?, ?, ?, 'pending', ?)",
		uid, "test", "websocket", title, body, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "Test notification queued"})
}

// --- Vault search (§17) ---
func handleVaultSearch(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		TermHashes []string `json:"term_hashes"`
		FolderID   *int64   `json:"folder_id,omitempty"`
		Limit      int      `json:"limit,omitempty"`
		Offset     int      `json:"offset,omitempty"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.TermHashes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "term_hashes array required"})
		return
	}
	if body.Limit <= 0 {
		body.Limit = 20
	}
	if body.Limit > 100 {
		body.Limit = 100
	}
	// Query: distinct file_id from vault_search_index where user_id and term_hash in (...)
	// Join vault_files to get file list, filter by folder_id if set
	q := `
		SELECT DISTINCT f.id, f.name, f.size_bytes, f.mime_type, f.folder_id, f.created_at, f.updated_at, f.storage_path
		FROM vault_files f
		INNER JOIN vault_search_index si ON si.file_id = f.id AND si.user_id = ?
		WHERE f.user_id = ? AND si.term_hash IN (?`
	args := []interface{}{uid, uid}
	placeholders := ""
	for i, h := range body.TermHashes {
		if i > 0 {
			placeholders += ",?"
		} else {
			placeholders = "?"
		}
		args = append(args, h)
	}
	q += placeholders + ")"
	if body.FolderID != nil {
		if *body.FolderID == 0 {
			q += " AND (f.folder_id IS NULL OR f.folder_id = 0)"
		} else {
			q += " AND f.folder_id = ?"
			args = append(args, *body.FolderID)
		}
	}
	q += " ORDER BY f.updated_at DESC LIMIT ? OFFSET ?"
	args = append(args, body.Limit, body.Offset)
	rows, err := db.DB.Query(q, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	defer rows.Close()
	var files []gin.H
	for rows.Next() {
		var id int64
		var name, mimeType, storagePath string
		var sizeBytes int64
		var folderID sql.NullInt64
		var createdAt, updatedAt int64
		if rows.Scan(&id, &name, &sizeBytes, &mimeType, &folderID, &createdAt, &updatedAt, &storagePath) != nil {
			continue
		}
		fid := interface{}(nil)
		if folderID.Valid {
			fid = folderID.Int64
		}
		files = append(files, gin.H{
			"id": id, "name": name, "size_bytes": sizeBytes, "mime_type": mimeType,
			"folder_id": fid, "created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if files == nil {
		files = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

// --- Admin (§18) ---
func handleAdminStats(c *gin.Context) {
	var users, products, orders int
	db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&users)
	db.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&products)
	db.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&orders)
	var reportsPending int
	db.DB.QueryRow("SELECT COUNT(*) FROM admin_reports WHERE status = 'pending'").Scan(&reportsPending)
	c.JSON(http.StatusOK, gin.H{
		"users":          users,
		"products":       products,
		"orders":         orders,
		"reports_pending": reportsPending,
	})
}

func handleAdminReportsList(c *gin.Context) {
	status := c.Query("status")
	q := "SELECT id, reporter_id, reported_type, reported_id, reason, status, created_at FROM admin_reports WHERE 1=1"
	args := []interface{}{}
	if status != "" {
		q += " AND status = ?"
		args = append(args, status)
	}
	q += " ORDER BY created_at DESC LIMIT 100"
	rows, err := db.DB.Query(q, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id, reporterID int64
		var reportedType, reportedID, reason, st string
		var createdAt int64
		if rows.Scan(&id, &reporterID, &reportedType, &reportedID, &reason, &st, &createdAt) != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "reporter_id": reporterID, "reported_type": reportedType,
			"reported_id": reportedID, "reason": reason, "status": st, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"reports": list})
}

func handleAdminReportGet(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var reporterID int64
	var reportedType, reportedID, reason, description, status, resolution string
	var assignedTo sql.NullInt64
	var createdAt, resolvedAt sql.NullInt64
	err = db.DB.QueryRow(
		"SELECT reporter_id, reported_type, reported_id, reason, description, status, assigned_to, resolution, created_at, resolved_at FROM admin_reports WHERE id = ?",
		id,
	).Scan(&reporterID, &reportedType, &reportedID, &reason, &description, &status, &assignedTo, &resolution, &createdAt, &resolvedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id, "reporter_id": reporterID, "reported_type": reportedType, "reported_id": reportedID,
		"reason": reason, "description": description, "status": status,
		"assigned_to": assignedTo.Int64, "resolution": resolution,
		"created_at": createdAt.Int64, "resolved_at": resolvedAt.Int64,
	})
}

func handleAdminReportAssign(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		AssignedTo int64 `json:"assigned_to"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assigned_to required"})
		return
	}
	res, err := db.DB.Exec("UPDATE admin_reports SET assigned_to = ? WHERE id = ?", body.AssignedTo, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleAdminReportResolve(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Resolution string `json:"resolution"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resolution required"})
		return
	}
	if body.Status == "" {
		body.Status = "resolved"
	}
	now := time.Now().Unix()
	res, err := db.DB.Exec(
		"UPDATE admin_reports SET status = ?, resolution = ?, resolved_at = ? WHERE id = ?",
		body.Status, body.Resolution, now, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleAdminUserGet(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var email, role string
	var createdAt sql.NullInt64
	err = db.DB.QueryRow(
		"SELECT email, role, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&email, &role, &createdAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var banned interface{}
	var banExpiresAt sql.NullInt64
	_ = db.DB.QueryRow(
		"SELECT expires_at FROM admin_bans WHERE user_id = ? AND lifted_at IS NULL ORDER BY created_at DESC LIMIT 1",
		userID,
	).Scan(&banExpiresAt)
	if banExpiresAt.Valid {
		banned = gin.H{"expires_at": banExpiresAt.Int64}
	} else {
		banned = nil
	}
	c.JSON(http.StatusOK, gin.H{
		"id": userID, "email": email, "role": role, "created_at": createdAt.Int64, "banned": banned,
	})
}

func handleAdminUserBan(c *gin.Context) {
	adminID := getUserID(c)
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Reason    string `json:"reason"`
		ExpiresAt *int64 `json:"expires_at,omitempty"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}
	var expiresAt interface{}
	if body.ExpiresAt != nil {
		expiresAt = *body.ExpiresAt
	}
	now := time.Now().Unix()
	_, err = db.DB.Exec(
		"INSERT INTO admin_bans (user_id, banned_by, reason, expires_at, created_at) VALUES (?, ?, ?, ?, ?)",
		userID, adminID, body.Reason, expiresAt, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleAdminUserUnban(c *gin.Context) {
	adminID := getUserID(c)
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	res, err := db.DB.Exec(
		"UPDATE admin_bans SET lifted_at = ?, lifted_by = ? WHERE user_id = ? AND lifted_at IS NULL",
		time.Now().Unix(), adminID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active ban"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func handleReportCreate(c *gin.Context) {
	uid := getUserID(c)
	var body struct {
		ReportedType string `json:"reported_type"`
		ReportedID   string `json:"reported_id"`
		Reason       string `json:"reason"`
		Description  string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ReportedType == "" || body.ReportedID == "" || body.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reported_type, reported_id, reason required"})
		return
	}
	_, err := db.DB.Exec(
		"INSERT INTO admin_reports (reporter_id, reported_type, reported_id, reason, description, status) VALUES (?, ?, ?, ?, ?, 'pending')",
		uid, body.ReportedType, body.ReportedID, body.Reason, body.Description,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
