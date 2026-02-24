package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"omnixius-api/db"

	"github.com/gin-gonic/gin"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	if err := db.Open(":memory:"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.DB.Close() })
	cfg = LoadConfig()
	gin.SetMode(gin.TestMode)
}

func TestLogin_UnknownEmail_Returns401(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"nobody@example.com","password":"any"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	handleLogin(c)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", w.Code)
	}
}

func TestRegister_Valid_Returns201AndToken(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"email":"u@test.com","password":"password123","name":"Test"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	handleRegister(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("got status %d, want 201", w.Code)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	user, _ := out["user"].(map[string]interface{})
	if user == nil {
		t.Fatal("expected user object in response")
	}
	if user["email"] != "u@test.com" {
		t.Errorf("user.email: got %v", user["email"])
	}
	if tok, _ := out["token"].(string); tok == "" {
		t.Error("expected non-empty token in response")
	}
}

func TestProductGet_NotFound_Returns404(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/products/99999", nil)
	c.Params = gin.Params{{Key: "id", Value: "99999"}}
	handleProductGet(c)
	if w.Code != http.StatusNotFound {
		t.Errorf("got status %d, want 404", w.Code)
	}
}

func TestOrderCreate_NoAuth_Returns401(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(`{"product_id":1,"quantity":1}`))
	c.Request.Header.Set("Content-Type", "application/json")
	handleOrderCreate(c)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", w.Code)
	}
}

func TestMessageSend_NoAuth_Returns401(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/messages/conversation/1", strings.NewReader(`{"body":"hello"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handleMessageSend(c)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", w.Code)
	}
}
