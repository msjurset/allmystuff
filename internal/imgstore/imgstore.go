package imgstore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type ImageStore struct {
	baseDir string
}

func New(baseDir string) *ImageStore {
	return &ImageStore{baseDir: baseDir}
}

func (s *ImageStore) Save(itemID, imageID uuid.UUID, ext string, r io.Reader) (string, error) {
	dir := filepath.Join(s.baseDir, itemID.String())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create image dir: %w", err)
	}

	filename := imageID.String() + ext
	relPath := filepath.Join(itemID.String(), filename)
	fullPath := filepath.Join(s.baseDir, relPath)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create image file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("write image file: %w", err)
	}

	return relPath, nil
}

func (s *ImageStore) Path(relPath string) string {
	return filepath.Join(s.baseDir, relPath)
}

func (s *ImageStore) Delete(relPath string) error {
	fullPath := filepath.Join(s.baseDir, relPath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete image file: %w", err)
	}
	return nil
}

func (s *ImageStore) DeleteItemDir(itemID uuid.UUID) error {
	dir := filepath.Join(s.baseDir, itemID.String())
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("delete item image dir: %w", err)
	}
	return nil
}
