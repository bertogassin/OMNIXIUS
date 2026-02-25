// Product service layer: get and create. Handlers do validation and file upload, then call service.
package main

import (
	"database/sql"
	"errors"
	"strconv"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

var ErrProductNotFound = errors.New("product not found")

type productRow struct {
	ID             int64
	UserID         int64
	Title          string
	Description    string
	Price          float64
	Category       string
	Location       string
	ImagePath      sql.NullString
	CreatedAt      int64
	IsService      int
	IsSubscription int
	SellerName     sql.NullString
	SellerEmail    string
}

func (p productRow) toH() gin.H {
	return gin.H{
		"id": p.ID, "user_id": p.UserID, "title": p.Title, "description": p.Description,
		"price": p.Price, "category": p.Category, "location": p.Location,
		"image_path": p.ImagePath.String, "created_at": p.CreatedAt, "is_service": p.IsService, "is_subscription": p.IsSubscription,
		"seller_id": p.UserID, "seller_name": p.SellerName.String, "seller_email": p.SellerEmail,
	}
}

// ProductGet returns a product by ID for public view (with seller name/email, seller_verified). Returns ErrProductNotFound if missing.
func ProductGet(idStr string) (gin.H, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, ErrProductNotFound
	}
	var p productRow
	var emailVerified, phoneVerified int
	err = db.DB.QueryRow(
		`SELECT p.id, p.user_id, p.title, p.description, p.price, p.category, p.location, p.image_path, p.created_at, COALESCE(p.is_service, 0), COALESCE(p.is_subscription, 0), u.name, u.email, COALESCE(u.email_verified, 0), COALESCE(u.phone_verified, 0) FROM products p JOIN users u ON u.id = p.user_id WHERE p.id = ?`,
		id,
	).Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &p.Price, &p.Category, &p.Location, &p.ImagePath, &p.CreatedAt, &p.IsService, &p.IsSubscription, &p.SellerName, &p.SellerEmail, &emailVerified, &phoneVerified)
	if err != nil {
		return nil, ErrProductNotFound
	}
	h := p.toH()
	h["seller_verified"] = emailVerified == 1 || phoneVerified == 1
	return h, nil
}

// ProductClosedContent returns closed_content_url, ownerID, is_subscription for product. Caller checks subscription/owner.
func ProductClosedContent(idStr string) (url string, ownerID int64, isSub int, err error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return "", 0, 0, ErrProductNotFound
	}
	var u sql.NullString
	if db.DB.QueryRow("SELECT closed_content_url, user_id, COALESCE(is_subscription, 0) FROM products WHERE id = ?", id).Scan(&u, &ownerID, &isSub) != nil {
		return "", 0, 0, ErrProductNotFound
	}
	return u.String, ownerID, isSub, nil
}

// ProductCreate inserts a product and returns its public map. isService, isSubscription: 0 or 1. closedContentURL optional.
func ProductCreate(ownerID int64, title, description, category, location, imagePath string, price float64, isService, isSubscription int, closedContentURL string) (gin.H, error) {
	if isService != 1 {
		isService = 0
	}
	if isSubscription != 1 {
		isSubscription = 0
	}
	if len(closedContentURL) > 2048 {
		closedContentURL = closedContentURL[:2048]
	}
	res, err := db.DB.Exec(
		"INSERT INTO products (user_id, title, description, price, category, location, image_path, is_service, is_subscription, closed_content_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		ownerID, title, description, price, category, nullStr(location), nullStr(imagePath), isService, isSubscription, nullStr(closedContentURL),
	)
	if err != nil {
		return nil, err
	}
	newID, _ := res.LastInsertId()
	var p productRow
	err = db.DB.QueryRow("SELECT id, user_id, title, description, price, category, location, image_path, created_at, COALESCE(is_service, 0), COALESCE(is_subscription, 0) FROM products WHERE id = ?", newID).
		Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &p.Price, &p.Category, &p.Location, &p.ImagePath, &p.CreatedAt, &p.IsService, &p.IsSubscription)
	if err != nil {
		return nil, err
	}
	return p.toH(), nil
}
