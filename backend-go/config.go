package main

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"omnixius-api/pqc"
)

type Config struct {
	Port             string
	SiteRoot         string // directory with index.html, app/, css/ (default: parent of backend-go)
	AppURL           string // frontend app base URL for redirect after register (e.g. https://bertogassin.github.io/OMNIXIUS)
	AllowedOrigins   string // comma-separated; empty = "*" (dev)
	MaxLoginAttempts int
	UploadDir        string
	MaxFileSize      int64
	// Argon2id (quantum-resistant KDF)
	Argon2Time    uint32
	Argon2Memory  uint32
	Argon2Threads uint8
	// Dilithium3 (post-quantum auth tokens)
	PQCPublicKey  []byte
	PQCPrivateKey []byte
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	mem := uint32(64 * 1024) // 64 MiB
	if v := os.Getenv("ARGON2_MEMORY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			mem = uint32(n)
		}
	}
	appURL := strings.TrimSuffix(os.Getenv("APP_URL"), "/")
	siteRoot := os.Getenv("SITE_ROOT")
	if siteRoot == "" {
		siteRoot = ".." // from backend-go/ run dir = project root
	}
	siteRoot = filepath.Clean(siteRoot)
	cfg := Config{
		Port:             port,
		SiteRoot:         siteRoot,
		AppURL:           appURL,
		AllowedOrigins:   os.Getenv("ALLOWED_ORIGINS"), // e.g. "https://bertogassin.github.io,https://omnixius.com"
		MaxLoginAttempts: 5,
		UploadDir:        "uploads",
		MaxFileSize:      5 * 1024 * 1024,
		Argon2Time:      3,
		Argon2Memory:    mem,
		Argon2Threads:    2,
	}
	// PQC keys from env (base64). Required for production.
	if b, err := base64.StdEncoding.DecodeString(os.Getenv("DILITHIUM_PUBLIC_KEY")); err == nil && len(b) > 0 {
		cfg.PQCPublicKey = b
	}
	if b, err := base64.StdEncoding.DecodeString(os.Getenv("DILITHIUM_PRIVATE_KEY")); err == nil && len(b) > 0 {
		cfg.PQCPrivateKey = b
	}
	if len(cfg.PQCPublicKey) == 0 || len(cfg.PQCPrivateKey) == 0 {
		priv, pub, err := pqc.GenerateKey()
		if err != nil {
			log.Fatal("pqc.GenerateKey: ", err)
		}
		cfg.PQCPrivateKey = priv
		cfg.PQCPublicKey = pub
		log.Print("PQC: DILITHIUM_PUBLIC_KEY/DILITHIUM_PRIVATE_KEY not set; ephemeral keys generated (tokens invalid after restart). Set env in production.")
	}
	return cfg
}
