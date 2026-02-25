// Balance service: user balance (internal units). Stub for Trade — subscriptions, rewards, payouts later.
package main

import (
	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

// BalanceGet returns current user balance (0 if no row).
func BalanceGet(userID int64) gin.H {
	var balance float64
	err := db.DB.QueryRow("SELECT balance FROM user_balances WHERE user_id = ?", userID).Scan(&balance)
	if err != nil {
		balance = 0
	}
	return gin.H{"balance": balance}
}

// BalanceCredit adds amount to user balance (stub: test credit). Returns new balance.
func BalanceCredit(userID int64, amount float64) (gin.H, error) {
	if amount <= 0 {
		return nil, nil // no-op
	}
	_, err := db.DB.Exec(
		`INSERT INTO user_balances (user_id, balance, updated_at) VALUES (?, ?, unixepoch())
		 ON CONFLICT(user_id) DO UPDATE SET balance = balance + excluded.balance, updated_at = unixepoch()`,
		userID, amount,
	)
	if err != nil {
		return nil, err
	}
	var newBalance float64
	db.DB.QueryRow("SELECT balance FROM user_balances WHERE user_id = ?", userID).Scan(&newBalance)
	return gin.H{"balance": newBalance, "credited": amount}, nil
}
