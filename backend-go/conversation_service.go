// Conversation and message service layer. Handlers validate and call service.
package main

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

var (
	ErrConvUserNotFound   = errors.New("user not found")
	ErrConvForbidden     = errors.New("forbidden")
	ErrMessageNotFound   = errors.New("message not found")
)

// ConversationsList returns conversations for the user (batched: other user + last message per conv).
func ConversationsList(uid int64) []gin.H {
	rows, _ := db.DB.Query(`SELECT c.id, c.product_id, c.updated_at FROM conversations c JOIN conversation_participants cp ON cp.conversation_id = c.id WHERE cp.user_id = ? ORDER BY c.updated_at DESC`, uid)
	type convRow struct{ id, productID, updated int64 }
	var convs []convRow
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var cid, pid, updated int64
			var pidNull sql.NullInt64
			rows.Scan(&cid, &pidNull, &updated)
			if pidNull.Valid {
				pid = pidNull.Int64
			}
			convs = append(convs, convRow{cid, pid, updated})
		}
	}
	if len(convs) == 0 {
		return []gin.H{}
	}
	cids := make([]interface{}, 0, len(convs))
	for _, c := range convs {
		cids = append(cids, c.id)
	}
	placeholders := strings.Repeat("?,", len(cids))
	placeholders = placeholders[:len(placeholders)-1]
	qry := `SELECT cp.conversation_id, u.id, u.name, u.email FROM conversation_participants cp JOIN users u ON u.id = cp.user_id WHERE cp.conversation_id IN (` + placeholders + `) AND cp.user_id != ?`
	args := append(cids, uid)
	otherRows, _ := db.DB.Query(qry, args...)
	otherByConv := make(map[int64]gin.H)
	if otherRows != nil {
		defer otherRows.Close()
		for otherRows.Next() {
			var cid, otherID int64
			var name, email sql.NullString
			otherRows.Scan(&cid, &otherID, &name, &email)
			otherByConv[cid] = gin.H{"id": otherID, "name": name.String, "email": email.String}
		}
	}
	lastRows, _ := db.DB.Query(`SELECT conversation_id, body FROM (SELECT conversation_id, body, ROW_NUMBER() OVER (PARTITION BY conversation_id ORDER BY created_at DESC) as rn FROM messages) WHERE rn = 1`)
	lastByConv := make(map[int64]string)
	if lastRows != nil {
		defer lastRows.Close()
		for lastRows.Next() {
			var cid int64
			var body sql.NullString
			lastRows.Scan(&cid, &body)
			lastByConv[cid] = body.String
		}
	}
	// Unread per conversation: messages where sender != uid and read_at IS NULL
	unreadRows, _ := db.DB.Query(`SELECT conversation_id, COUNT(*) FROM messages WHERE conversation_id IN (`+placeholders+`) AND sender_id != ? AND read_at IS NULL GROUP BY conversation_id`, append(cids, uid)...)
	unreadByConv := make(map[int64]int64)
	if unreadRows != nil {
		defer unreadRows.Close()
		for unreadRows.Next() {
			var cid int64
			var n int64
			unreadRows.Scan(&cid, &n)
			unreadByConv[cid] = n
		}
	}
	list := make([]gin.H, 0, len(convs))
	for _, c := range convs {
		other := otherByConv[c.id]
		if other == nil {
			other = gin.H{"id": int64(0), "name": "", "email": ""}
		}
		unread := unreadByConv[c.id]
		list = append(list, gin.H{"id": c.id, "product_id": c.productID, "updated_at": c.updated, "last_message": lastByConv[c.id], "other": other, "unread": unread > 0})
	}
	return list
}

// UnreadCount returns total count of unread messages for the user (messages from others, not yet read).
func UnreadCount(uid int64) int64 {
	var n int64
	db.DB.QueryRow(`
		SELECT COUNT(*) FROM messages m
		JOIN conversation_participants cp ON cp.conversation_id = m.conversation_id AND cp.user_id = ?
		WHERE m.sender_id != ? AND m.read_at IS NULL
	`, uid, uid).Scan(&n)
	return n
}

// ConversationCreate finds or creates a conversation between uid and otherUserID, optional productID. Returns conversation ID.
func ConversationCreate(uid, otherUserID, productID int64) (convID int64, err error) {
	if otherUserID == 0 || otherUserID == uid {
		return 0, ErrConvUserNotFound
	}
	var ok int64
	if db.DB.QueryRow("SELECT id FROM users WHERE id = ?", otherUserID).Scan(&ok) != nil {
		return 0, ErrConvUserNotFound
	}
	err = db.DB.QueryRow(
		`SELECT c.id FROM conversations c JOIN conversation_participants cp1 ON cp1.conversation_id = c.id AND cp1.user_id = ? JOIN conversation_participants cp2 ON cp2.conversation_id = c.id AND cp2.user_id = ? WHERE (c.product_id IS NULL AND ? = 0) OR c.product_id = ?`,
		uid, otherUserID, productID, productID,
	).Scan(&convID)
	if err != nil {
		res, _ := db.DB.Exec("INSERT INTO conversations (product_id) VALUES (?)", nullInt64(productID))
		convID, _ = res.LastInsertId()
		db.DB.Exec("INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?), (?, ?)", convID, uid, convID, otherUserID)
	}
	return convID, nil
}

// MessagesList returns messages in a conversation; returns error if user is not a participant.
func MessagesList(convID int64, uid int64) ([]gin.H, error) {
	var ok int
	if db.DB.QueryRow("SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?", convID, uid).Scan(&ok) != nil {
		return nil, ErrConvForbidden
	}
	rows, _ := db.DB.Query("SELECT m.id, m.sender_id, m.body, m.read_at, m.created_at, u.name FROM messages m JOIN users u ON u.id = m.sender_id WHERE m.conversation_id = ? ORDER BY m.created_at ASC", convID)
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
	return list, nil
}

// MessageSend adds a message; returns created message map or error if not a participant.
func MessageSend(convID int64, uid int64, body string) (gin.H, error) {
	var ok int
	if db.DB.QueryRow("SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?", convID, uid).Scan(&ok) != nil {
		return nil, ErrConvForbidden
	}
	res, err := db.DB.Exec("INSERT INTO messages (conversation_id, sender_id, body) VALUES (?, ?, ?)", convID, uid, body)
	if err != nil {
		return nil, err
	}
	db.DB.Exec("UPDATE conversations SET updated_at = unixepoch() WHERE id = ?", convID)
	mid, _ := res.LastInsertId()
	return gin.H{"id": mid, "conversation_id": convID, "sender_id": uid, "body": body, "created_at": time.Now().Unix()}, nil
}

// MessageMarkRead marks a message as read; returns error if not participant or message not found.
func MessageMarkRead(messageID int64, uid int64) error {
	var senderID, convID int64
	if db.DB.QueryRow("SELECT sender_id, conversation_id FROM messages WHERE id = ?", messageID).Scan(&senderID, &convID) != nil {
		return ErrMessageNotFound
	}
	if senderID == uid {
		return nil
	}
	var ok int
	if db.DB.QueryRow("SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?", convID, uid).Scan(&ok) != nil {
		return ErrConvForbidden
	}
	db.DB.Exec("UPDATE messages SET read_at = unixepoch() WHERE id = ? AND read_at IS NULL", messageID)
	return nil
}

