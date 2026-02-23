package db

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

var DB *sql.DB

func Open(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	DB.Exec("PRAGMA foreign_keys = ON")
	schema, _ := schemaFS.ReadFile("schema.sql")
	_, err = DB.Exec(string(schema))
	return err
}

func InitUploadDirs(uploadDir string) {
	dirs := []string{uploadDir, filepath.Join(uploadDir, "products"), filepath.Join(uploadDir, "avatars")}
	for _, d := range dirs {
		os.MkdirAll(d, 0755)
	}
}
