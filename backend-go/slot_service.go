// Slot service: list/add slots for service products; book slot → order + mail notification.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

var (
	ErrSlotProductNotFound = errors.New("product not found")
	ErrSlotForbidden       = errors.New("not the product owner")
	ErrSlotNotFound        = errors.New("slot not found")
	ErrSlotNotFree         = errors.New("slot already booked")
	ErrSlotNotService      = errors.New("product is not a service")
)

// SlotsList returns slots for a product. For owner: all slots; for others: only free.
func SlotsList(productID string, viewerID int64) ([]gin.H, error) {
	pid, err := strconv.ParseInt(productID, 10, 64)
	if err != nil || pid <= 0 {
		return nil, ErrSlotProductNotFound
	}
	var ownerID int64
	var isService int
	if db.DB.QueryRow("SELECT user_id, COALESCE(is_service, 0) FROM products WHERE id = ?", pid).Scan(&ownerID, &isService) != nil {
		return nil, ErrSlotProductNotFound
	}
	qry := `SELECT id, product_id, slot_at, status, order_id, created_at FROM product_slots WHERE product_id = ?`
	args := []interface{}{pid}
	if viewerID != ownerID {
		qry += ` AND status = 'free'`
	}
	qry += ` ORDER BY slot_at ASC`
	rows, err := db.DB.Query(qry, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id, slotAt, created int64
		var orderID sql.NullInt64
		var productID int64
		var status string
		rows.Scan(&id, &productID, &slotAt, &status, &orderID, &created)
		list = append(list, gin.H{"id": id, "product_id": productID, "slot_at": slotAt, "status": status, "order_id": orderID.Int64, "created_at": created})
	}
	return list, nil
}

// SlotsAdd adds a slot (seller only). slotAt is Unix timestamp.
func SlotsAdd(productID string, ownerID int64, slotAt int64) (gin.H, error) {
	pid, err := strconv.ParseInt(productID, 10, 64)
	if err != nil || pid <= 0 {
		return nil, ErrSlotProductNotFound
	}
	var dbOwnerID int64
	if db.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", pid).Scan(&dbOwnerID) != nil {
		return nil, ErrSlotProductNotFound
	}
	if dbOwnerID != ownerID {
		return nil, ErrSlotForbidden
	}
	if slotAt <= 0 {
		slotAt = time.Now().Unix()
	}
	res, err := db.DB.Exec("INSERT INTO product_slots (product_id, slot_at, status) VALUES (?, ?, 'free')", pid, slotAt)
	if err != nil {
		return nil, err
	}
	sid, _ := res.LastInsertId()
	return gin.H{"id": sid, "product_id": pid, "slot_at": slotAt, "status": "free"}, nil
}

// SlotBook books a slot: creates order, marks slot booked, notifies seller in Mail.
func SlotBook(productID, slotID string, buyerID int64) (gin.H, error) {
	pid, err := strconv.ParseInt(productID, 10, 64)
	if err != nil || pid <= 0 {
		return nil, ErrSlotProductNotFound
	}
	sid, err := strconv.ParseInt(slotID, 10, 64)
	if err != nil || sid <= 0 {
		return nil, ErrSlotNotFound
	}
	var sellerID int64
	var isService int
	if db.DB.QueryRow("SELECT user_id, COALESCE(is_service, 0) FROM products WHERE id = ?", pid).Scan(&sellerID, &isService) != nil {
		return nil, ErrSlotProductNotFound
	}
	if isService != 1 {
		return nil, ErrSlotNotService
	}
	if sellerID == buyerID {
		return nil, ErrOrderOwnProduct
	}
	var status string
	var slotProductID int64
	if db.DB.QueryRow("SELECT status, product_id FROM product_slots WHERE id = ?", sid).Scan(&status, &slotProductID) != nil {
		return nil, ErrSlotNotFound
	}
	if slotProductID != pid {
		return nil, ErrSlotNotFound
	}
	if status != "free" {
		return nil, ErrSlotNotFree
	}
	// Create order with slot_id
	order, err := OrderCreateWithSlot(buyerID, pid, sid, "", false)
	if err != nil {
		return nil, err
	}
	oid := order["id"].(int64)
	db.DB.Exec("UPDATE product_slots SET status = 'booked', order_id = ? WHERE id = ?", oid, sid)
	// Notify seller: get or create conversation, send message
	convID, _ := ConversationCreate(buyerID, sellerID, pid)
	var slotAtUnix int64
	db.DB.QueryRow("SELECT slot_at FROM product_slots WHERE id = ?", sid).Scan(&slotAtUnix)
	slotTime := time.Unix(slotAtUnix, 0).Format("2006-01-02 15:04")
	msg := fmt.Sprintf("Booked your service for %s. Order #%d.", slotTime, oid)
	MessageSend(convID, buyerID, msg)
	return gin.H{"order": order, "slot_id": sid, "message": "Booked; seller notified."}, nil
}
