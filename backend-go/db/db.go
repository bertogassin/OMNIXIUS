package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

//go:embed migrations/*.sql
var migrationsFS embed.FS

var DB *sql.DB

func Open(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	DB.Exec("PRAGMA foreign_keys = ON")
	schema, _ := schemaFS.ReadFile("schema.sql")
	if _, err = DB.Exec(string(schema)); err != nil {
		return err
	}
	return RunMigrations()
}

// RunMigrations runs embedded migrations in order (002_*.sql, 003_*.sql, ...).
func RunMigrations() error {
	var version int
	if err := DB.QueryRow("SELECT version FROM schema_version WHERE id = 1").Scan(&version); err != nil {
		return fmt.Errorf("schema_version: %w", err)
	}
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	var files []struct {
		n    int
		name string
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		n, err := strconv.Atoi(e.Name()[:3])
		if err != nil || n < 2 {
			continue
		}
		files = append(files, struct{ n int; name string }{n: n, name: e.Name()})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].n < files[j].n })
	for _, f := range files {
		if f.n <= version {
			continue
		}
		body, err := migrationsFS.ReadFile("migrations/" + f.name)
		if err != nil {
			return fmt.Errorf("migration %s: %w", f.name, err)
		}
		sql := strings.TrimSpace(string(body))
		if sql != "" && !strings.HasPrefix(sql, "--") {
			if _, err = DB.Exec(sql); err != nil {
				return fmt.Errorf("migration %s: %w", f.name, err)
			}
		}
		if _, err = DB.Exec("UPDATE schema_version SET version = ? WHERE id = 1", f.n); err != nil {
			return err
		}
		version = f.n
	}
	return nil
}

func InitUploadDirs(uploadDir string) {
	dirs := []string{uploadDir, filepath.Join(uploadDir, "products"), filepath.Join(uploadDir, "avatars"), filepath.Join(uploadDir, "vault")}
	for _, d := range dirs {
		os.MkdirAll(d, 0755)
	}
}
