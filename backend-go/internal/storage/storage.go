// Package storage provides §1.5 StorageProvider and Local implementation (doc v4.0).
package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

// StorageItem holds metadata for a stored object (§1.5.1).
type StorageItem struct {
	Path         string
	Size         int64
	ETag         string
	ContentType  string
	LastModified time.Time
	Metadata     map[string]string
}

// StorageProvider defines storage operations (§1.5.1).
type StorageProvider interface {
	Put(path string, data []byte, metadata map[string]string) error
	Get(path string) ([]byte, error)
	Delete(path string) error
	List(prefix string) ([]StorageItem, error)
	Head(path string) (StorageItem, error)
	GenerateUploadURL(path string, expires time.Duration) (string, error)
	GenerateDownloadURL(path string, expires time.Duration) (string, error)
}

// LocalStorage stores files on disk. Pre-signed URLs return error (use Put/Get).
type LocalStorage struct {
	BaseDir string
}

// NewLocalStorage returns a LocalStorage.
func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir}
}

func (s *LocalStorage) fullPath(path string) string {
	return filepath.Join(s.BaseDir, filepath.Clean(path))
}

func (s *LocalStorage) Put(path string, data []byte, metadata map[string]string) error {
	full := s.fullPath(path)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	return os.WriteFile(full, data, 0644)
}

func (s *LocalStorage) Get(path string) ([]byte, error) {
	return os.ReadFile(s.fullPath(path))
}

func (s *LocalStorage) Delete(path string) error {
	return os.Remove(s.fullPath(path))
}

func (s *LocalStorage) List(prefix string) ([]StorageItem, error) {
	dir := s.fullPath(prefix)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []StorageItem
	for _, e := range entries {
		full := filepath.Join(dir, e.Name())
		rel, _ := filepath.Rel(s.BaseDir, full)
		if e.IsDir() {
			out = append(out, StorageItem{Path: rel + "/", Size: 0})
			continue
		}
		info, _ := e.Info()
		out = append(out, StorageItem{Path: rel, Size: info.Size(), LastModified: info.ModTime()})
	}
	return out, nil
}

func (s *LocalStorage) Head(path string) (StorageItem, error) {
	info, err := os.Stat(s.fullPath(path))
	if err != nil {
		return StorageItem{}, err
	}
	rel, _ := filepath.Rel(s.BaseDir, s.fullPath(path))
	return StorageItem{Path: rel, Size: info.Size(), LastModified: info.ModTime()}, nil
}

func (s *LocalStorage) GenerateUploadURL(path string, expires time.Duration) (string, error) {
	return "", errors.New("pre-signed upload not supported for local storage")
}

func (s *LocalStorage) GenerateDownloadURL(path string, expires time.Duration) (string, error) {
	return "", errors.New("pre-signed download not supported for local storage")
}

// PutReader writes from reader to path.
func (s *LocalStorage) PutReader(path string, r io.Reader, metadata map[string]string) error {
	full := s.fullPath(path)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	f, err := os.Create(full)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
