// Remittance service: cross-border transfer requests (stub). Real transfer via Trade/IXI later.
package main

import (
	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

// RemittanceCreate records a remittance request. No actual transfer.
func RemittanceCreate(fromUserID int64, toIdentifier string, amount float64, currency string) (gin.H, error) {
	if currency == "" {
		currency = "USD"
	}
	res, err := db.DB.Exec(
		"INSERT INTO remittances (from_user_id, to_identifier, amount, currency, status) VALUES (?, ?, ?, ?, 'pending')",
		fromUserID, toIdentifier, amount, currency,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return gin.H{"id": id, "from_user_id": fromUserID, "to_identifier": toIdentifier, "amount": amount, "currency": currency, "status": "pending"}, nil
}

// RemittanceListMy returns remittance requests sent by the user, newest first.
func RemittanceListMy(fromUserID int64) ([]gin.H, error) {
	rows, err := db.DB.Query(
		"SELECT id, from_user_id, to_identifier, amount, currency, status, created_at FROM remittances WHERE from_user_id = ? ORDER BY created_at DESC",
		fromUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []gin.H
	for rows.Next() {
		var id, fromUserIDVal int64
		var toIdentifier, currency, status string
		var amount float64
		var createdAt *int64
		if err := rows.Scan(&id, &fromUserIDVal, &toIdentifier, &amount, &currency, &status, &createdAt); err != nil {
			return nil, err
		}
		h := gin.H{"id": id, "from_user_id": fromUserIDVal, "to_identifier": toIdentifier, "amount": amount, "currency": currency, "status": status}
		if createdAt != nil {
			h["created_at"] = *createdAt
		}
		list = append(list, h)
	}
	return list, rows.Err()
}
