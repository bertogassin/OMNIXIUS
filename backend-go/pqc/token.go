// Package pqc provides quantum-resistant auth tokens (Dilithium3 signatures).
package pqc

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"

	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

const (
	payloadLen     = 8 + 8       // userID (int64) + exp (int64) — legacy
	payloadLenV2   = 8 + 8 + 8   // userID + exp + sessionID (int64)
)

// GenerateKey creates a new Dilithium3 key pair. For production, load from env instead.
func GenerateKey() (privateKey, publicKey []byte, err error) {
	pk, sk, err := mode3.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return sk.Bytes(), pk.Bytes(), nil
}

// SignToken signs userID and expiry (exp) with the private key. Returns token string. Uses sessionID 0 (legacy).
func SignToken(privateKey []byte, userID int64, exp time.Time) (string, error) {
	return SignTokenWithSession(privateKey, userID, 0, exp)
}

// SignTokenWithSession signs userID, sessionID and expiry. SessionID 0 = legacy (no session row).
func SignTokenWithSession(privateKey []byte, userID int64, sessionID int64, exp time.Time) (string, error) {
	if len(privateKey) != mode3.PrivateKeySize {
		return "", errors.New("invalid private key size")
	}
	var sk mode3.PrivateKey
	if err := sk.UnmarshalBinary(privateKey); err != nil {
		return "", err
	}
	payload := make([]byte, payloadLenV2)
	binary.BigEndian.PutUint64(payload[0:8], uint64(userID))
	binary.BigEndian.PutUint64(payload[8:16], uint64(exp.Unix()))
	binary.BigEndian.PutUint64(payload[16:24], uint64(sessionID))
	sig := make([]byte, mode3.SignatureSize)
	mode3.SignTo(&sk, payload, sig)
	token := base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(sig)
	return token, nil
}

// VerifyToken verifies the token and returns userID, expiry and sessionID (0 if legacy 16-byte payload).
func VerifyToken(publicKey []byte, token string) (userID int64, exp time.Time, sessionID int64, err error) {
	userID, exp, sessionID, err = verifyTokenPayload(publicKey, token)
	return userID, exp, sessionID, err
}

func verifyTokenPayload(publicKey []byte, token string) (userID int64, exp time.Time, sessionID int64, err error) {
	if len(publicKey) != mode3.PublicKeySize {
		return 0, time.Time{}, 0, errors.New("invalid public key size")
	}
	var pk mode3.PublicKey
	if err := pk.UnmarshalBinary(publicKey); err != nil {
		return 0, time.Time{}, 0, err
	}
	i := 0
	for i < len(token) && token[i] != '.' {
		i++
	}
	if i == len(token) {
		return 0, time.Time{}, 0, errors.New("invalid token format")
	}
	payloadB64 := token[:i]
	sigB64 := token[i+1:]
	payload, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return 0, time.Time{}, 0, errors.New("invalid token payload")
	}
	if len(payload) != payloadLen && len(payload) != payloadLenV2 {
		return 0, time.Time{}, 0, errors.New("invalid token payload")
	}
	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil || len(sig) != mode3.SignatureSize {
		return 0, time.Time{}, 0, errors.New("invalid token signature")
	}
	if !mode3.Verify(&pk, payload, sig) {
		return 0, time.Time{}, 0, errors.New("invalid signature")
	}
	userID = int64(binary.BigEndian.Uint64(payload[0:8]))
	expUnix := int64(binary.BigEndian.Uint64(payload[8:16]))
	exp = time.Unix(expUnix, 0)
	if time.Now().After(exp) {
		return 0, time.Time{}, 0, errors.New("token expired")
	}
	if len(payload) >= payloadLenV2 {
		sessionID = int64(binary.BigEndian.Uint64(payload[16:24]))
	}
	return userID, exp, sessionID, nil
}
