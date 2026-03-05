package imgstore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	s := New("/tmp/test-images")
	if s == nil {
		t.Fatal("expected non-nil ImageStore")
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	itemID := uuid.New()
	imageID := uuid.New()
	content := "fake image data"

	relPath, err := s.Save(itemID, imageID, ".jpg", strings.NewReader(content))
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify relative path format
	expected := filepath.Join(itemID.String(), imageID.String()+".jpg")
	if relPath != expected {
		t.Errorf("expected relPath %q, got %q", expected, relPath)
	}

	// Verify file exists and has correct content
	data, err := os.ReadFile(filepath.Join(dir, relPath))
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	if string(data) != content {
		t.Errorf("expected content %q, got %q", content, string(data))
	}
}

func TestSave_CreatesItemDirectory(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	itemID := uuid.New()
	imageID := uuid.New()

	_, err := s.Save(itemID, imageID, ".png", strings.NewReader("data"))
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	itemDir := filepath.Join(dir, itemID.String())
	info, err := os.Stat(itemDir)
	if err != nil {
		t.Fatalf("item directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory, got file")
	}
}

func TestSave_MultipleImages(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	itemID := uuid.New()

	for i := 0; i < 3; i++ {
		imageID := uuid.New()
		_, err := s.Save(itemID, imageID, ".jpg", strings.NewReader("data"))
		if err != nil {
			t.Fatalf("Save #%d failed: %v", i, err)
		}
	}

	entries, err := os.ReadDir(filepath.Join(dir, itemID.String()))
	if err != nil {
		t.Fatalf("failed to read item dir: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 files, got %d", len(entries))
	}
}

func TestPath(t *testing.T) {
	s := New("/data/images")
	p := s.Path("abc/def.jpg")
	expected := filepath.Join("/data/images", "abc", "def.jpg")
	if p != expected {
		t.Errorf("expected %q, got %q", expected, p)
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	itemID := uuid.New()
	imageID := uuid.New()

	relPath, err := s.Save(itemID, imageID, ".jpg", strings.NewReader("data"))
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if err := s.Delete(relPath); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, relPath)); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestDelete_NonExistent(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	// Should not error on non-existent file
	err := s.Delete("nonexistent/file.jpg")
	if err != nil {
		t.Errorf("expected no error for non-existent file, got: %v", err)
	}
}

func TestDeleteItemDir(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	itemID := uuid.New()

	// Create a couple of images
	for i := 0; i < 2; i++ {
		_, err := s.Save(itemID, uuid.New(), ".jpg", strings.NewReader("data"))
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	if err := s.DeleteItemDir(itemID); err != nil {
		t.Fatalf("DeleteItemDir failed: %v", err)
	}

	itemDir := filepath.Join(dir, itemID.String())
	if _, err := os.Stat(itemDir); !os.IsNotExist(err) {
		t.Error("expected item directory to be removed")
	}
}

func TestDeleteItemDir_NonExistent(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	// Should not error on non-existent directory
	err := s.DeleteItemDir(uuid.New())
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
