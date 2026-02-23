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
	ID          int64
	UserID      int64
	Title       string
	Description string
	Price       float64
	Category    string
	Location    string
	ImagePath   sql.NullString
	CreatedAt   int64
	SellerName  sql.NullString
	SellerEmail string
}

func (p productRow) toH() gin.H {
	return gin.H{
		"id": p.ID, "user_id": p.UserID, "title": p.Title, "description": p.Description,
		"price": p.Price, "category": p.Category, "location": p.Location,
		"image_path": p.ImagePath.String, "created_at": p.CreatedAt,
		"seller_id": p.UserID, "seller_name": p.SellerName.String, "seller_email": p.SellerEmail,
	}
}

// ProductGet returns a product by ID for public view (with seller name/email). Returns ErrProductNotFound if missing.
func ProductGet(idStr string) (gin.H, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, ErrProductNotFound
	}
	var p productRow
	err = db.DB.QueryRow(
		`SELECT p.id, p.user_id, p.title, p.description, p.price, p.category, p.location, p.image_path, p.created_at, u.name, u.email FROM products p JOIN users u ON u.id = p.user_id WHERE p.id = ?`,
		id,
	).Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &p.Price, &p.Category, &p.Location, &p.ImagePath, &p.CreatedAt, &p.SellerName, &p.SellerEmail)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return p.toH(), nil
}

// ProductCreate inserts a product and returns its public map. Caller must validate inputs and optional image path.
func ProductCreate(ownerID int64, title, description, category, location, imagePath string, price float64) (gin.H, error) {
	res, err := db.DB.Exec(
		"INSERT INTO products (user_id, title, description, price, category, location, image_path) VALUES (?, ?, ?, ?, ?, ?, ?)",
		ownerID, title, description, price, category, nullStr(location), nullStr(imagePath),
	)
	if err != nil {
		return nil, err
	}
	newID, _ := res.LastInsertId()
	var p productRow
	err = db.DB.QueryRow("SELECT id, user_id, title, description, price, category, location, image_path, created_at FROM products WHERE id = ?", newID).
		Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &p.Price, &p.Category, &p.Location, &p.ImagePath, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p.toH(), nil
}
