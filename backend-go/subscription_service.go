// Subscription service: subscribe to a subscription-type product (stub; payments later).
package main

import (
	"database/sql"
	"errors"
	"strconv"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

var (
	ErrSubProductNotFound   = errors.New("product not found")
	ErrSubNotSubscription   = errors.New("product is not a subscription listing")
	ErrSubOwnProduct       = errors.New("cannot subscribe to own product")
	ErrSubAlreadySubscribed = errors.New("already subscribed")
)

// Subscribe creates a subscription (stub; no payment). Product must be is_subscription=1.
func Subscribe(productID string, userID int64) (gin.H, error) {
	pid, err := strconv.ParseInt(productID, 10, 64)
	if err != nil || pid <= 0 {
		return nil, ErrSubProductNotFound
	}
	var sellerID int64
	var isSub int
	if db.DB.QueryRow("SELECT user_id, COALESCE(is_subscription, 0) FROM products WHERE id = ?", pid).Scan(&sellerID, &isSub) != nil {
		return nil, ErrSubProductNotFound
	}
	if isSub != 1 {
		return nil, ErrSubNotSubscription
	}
	if sellerID == userID {
		return nil, ErrSubOwnProduct
	}
	var existing int64
	if db.DB.QueryRow("SELECT id FROM subscriptions WHERE product_id = ? AND user_id = ? AND status = 'active'", pid, userID).Scan(&existing) == nil {
		return nil, ErrSubAlreadySubscribed
	}
	res, err := db.DB.Exec("INSERT INTO subscriptions (product_id, user_id, status) VALUES (?, ?, 'active')", pid, userID)
	if err != nil {
		return nil, err
	}
	sid, _ := res.LastInsertId()
	return gin.H{"id": sid, "product_id": pid, "user_id": userID, "status": "active"}, nil
}

// SubscriptionsMy returns current user's active subscriptions with product details.
func SubscriptionsMy(userID int64) []gin.H {
	rows, _ := db.DB.Query(
		`SELECT s.id, s.product_id, s.status, s.created_at, p.title, p.price, p.image_path, u.id, u.name FROM subscriptions s JOIN products p ON p.id = s.product_id JOIN users u ON u.id = p.user_id WHERE s.user_id = ? AND s.status = 'active' ORDER BY s.created_at DESC`,
		userID,
	)
	var list []gin.H
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var sid, pid, created, sellerID int64
			var status, title string
			var price float64
			var imagePath, sellerName sql.NullString
			rows.Scan(&sid, &pid, &status, &created, &title, &price, &imagePath, &sellerID, &sellerName)
			list = append(list, gin.H{"id": sid, "product_id": pid, "status": status, "created_at": created, "title": title, "price": price, "image_path": imagePath.String, "seller_id": sellerID, "seller_name": sellerName.String})
		}
	}
	return list
}

// IsSubscribed returns true if user has an active subscription to the product.
func IsSubscribed(productID int64, userID int64) bool {
	if userID <= 0 {
		return false
	}
	var n int
	return db.DB.QueryRow("SELECT 1 FROM subscriptions WHERE product_id = ? AND user_id = ? AND status = 'active'", productID, userID).Scan(&n) == nil
}
