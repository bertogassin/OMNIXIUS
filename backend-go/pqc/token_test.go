package pqc

import (
	"testing"
	"time"
)

func TestGenerateKey(t *testing.T) {
	sk, pk, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	if len(sk) == 0 || len(pk) == 0 {
		t.Error("empty key returned")
	}
	// Sign/verify with same key pair
	tok, err := SignToken(sk, 1, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	uid, exp, _, err := VerifyToken(pk, tok)
	if err != nil {
		t.Fatal(err)
	}
	if uid != 1 {
		t.Errorf("userID: got %d want 1", uid)
	}
	if exp.Before(time.Now()) {
		t.Error("exp should be in future")
	}
}

func TestSignToken_InvalidKey(t *testing.T) {
	_, err := SignToken([]byte("short"), 1, time.Now().Add(time.Hour))
	if err == nil {
		t.Error("expected error for invalid key size")
	}
}

func TestVerifyToken_InvalidKey(t *testing.T) {
	_, _, _, err := VerifyToken([]byte("short"), "payload.sig")
	if err == nil {
		t.Error("expected error for invalid key size")
	}
}

func TestVerifyToken_InvalidFormat(t *testing.T) {
	_, pk, _ := GenerateKey()
	_, _, _, err := VerifyToken(pk, "no-dot")
	if err == nil {
		t.Error("expected error for invalid token format")
	}
}

func TestVerifyToken_Expired(t *testing.T) {
	sk, pk, _ := GenerateKey()
	tok, _ := SignToken(sk, 1, time.Now().Add(-time.Second))
	_, _, _, err := VerifyToken(pk, tok)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestVerifyToken_Tampered(t *testing.T) {
	sk, pk, _ := GenerateKey()
	tok, _ := SignToken(sk, 1, time.Now().Add(time.Hour))
	// Tamper: flip a byte in the payload part
	b := []byte(tok)
	for i := range b {
		if b[i] == '.' {
			b[i] = 'x'
			break
		}
	}
	_, _, _, err := VerifyToken(pk, string(b))
	if err == nil {
		t.Error("expected error for tampered token")
	}
}
