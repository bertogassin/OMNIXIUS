// Order service layer: list my orders, create order. Handlers validate input and call service.
package main

import (
	"database/sql"
	"errors"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

var (
	ErrOrderProductNotFound = errors.New("product not found")
	ErrOrderOwnProduct      = errors.New("cannot order own product")
)

// OrdersMy returns all orders where the user is buyer or seller.
func OrdersMy(userID int64) []gin.H {
	rows, _ := db.DB.Query(
		`SELECT o.id, o.product_id, o.buyer_id, o.seller_id, o.status, o.created_at, COALESCE(o.installment_plan, ''), p.title, p.price, p.image_path FROM orders o JOIN products p ON p.id = o.product_id WHERE o.buyer_id = ? OR o.seller_id = ? ORDER BY o.created_at DESC`,
		userID, userID,
	)
	var list []gin.H
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var oid, pid, buyer, seller, created int64
			var status, installmentPlan, title string
			var price float64
			var img sql.NullString
			rows.Scan(&oid, &pid, &buyer, &seller, &status, &created, &installmentPlan, &title, &price, &img)
			list = append(list, gin.H{"id": oid, "product_id": pid, "buyer_id": buyer, "seller_id": seller, "status": status, "created_at": created, "installment_plan": installmentPlan, "title": title, "price": price, "image_path": img.String})
		}
	}
	return list
}

// OrderCreate creates an order; buyer must not be the seller. installmentPlan: "" or "requested".
func OrderCreate(buyerID, productID int64, installmentPlan string) (gin.H, error) {
	return OrderCreateWithSlot(buyerID, productID, 0, installmentPlan)
}

// OrderCreateWithSlot creates an order, optionally linked to a slot (slotID 0 = none). installmentPlan: "" or "requested".
func OrderCreateWithSlot(buyerID, productID, slotID int64, installmentPlan string) (gin.H, error) {
	var sellerID int64
	if db.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", productID).Scan(&sellerID) != nil {
		return nil, ErrOrderProductNotFound
	}
	if sellerID == buyerID {
		return nil, ErrOrderOwnProduct
	}
	if installmentPlan != "requested" {
		installmentPlan = ""
	}
	var res sql.Result
	var err error
	if slotID > 0 {
		res, err = db.DB.Exec("INSERT INTO orders (product_id, buyer_id, seller_id, status, slot_id, installment_plan) VALUES (?, ?, ?, 'pending', ?, ?)", productID, buyerID, sellerID, slotID, installmentPlan)
	} else {
		res, err = db.DB.Exec("INSERT INTO orders (product_id, buyer_id, seller_id, status, installment_plan) VALUES (?, ?, ?, 'pending', ?)", productID, buyerID, sellerID, installmentPlan)
	}
	if err != nil {
		return nil, err
	}
	oid, _ := res.LastInsertId()
	out := gin.H{"id": oid, "product_id": productID, "buyer_id": buyerID, "seller_id": sellerID, "status": "pending", "installment_plan": installmentPlan}
	if slotID > 0 {
		out["slot_id"] = slotID
	}
	return out, nil
}
