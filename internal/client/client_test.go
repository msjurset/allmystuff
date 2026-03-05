package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"allmystuff/internal/model"
	"allmystuff/internal/store"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	c := New("http://localhost:8080", "mykey")
	if c.BaseURL != "http://localhost:8080" {
		t.Errorf("expected BaseURL http://localhost:8080, got %s", c.BaseURL)
	}
	if c.APIKey != "mykey" {
		t.Errorf("expected APIKey mykey, got %s", c.APIKey)
	}
	if c.HTTP == nil {
		t.Error("expected non-nil HTTP client")
	}
}

func TestClient_AuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]model.Item{})
	}))
	defer srv.Close()

	c := New(srv.URL, "test-key")
	_, _ = c.ListItems(store.ItemFilter{})

	if gotAuth != "Bearer test-key" {
		t.Errorf("expected 'Bearer test-key', got %q", gotAuth)
	}
}

func TestClient_NoAuthHeader_WhenKeyEmpty(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]model.Item{})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, _ = c.ListItems(store.ItemFilter{})

	if gotAuth != "" {
		t.Errorf("expected empty auth header, got %q", gotAuth)
	}
}

func TestClient_ListItems(t *testing.T) {
	items := []model.Item{
		{ID: uuid.New(), Name: "Test Item"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/items" {
			t.Errorf("expected path /api/items, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	result, err := c.ListItems(store.ItemFilter{})
	if err != nil {
		t.Fatalf("ListItems failed: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if result[0].Name != "Test Item" {
		t.Errorf("expected name 'Test Item', got %q", result[0].Name)
	}
}

func TestClient_ListItems_WithFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("q") != "camera" {
			t.Errorf("expected q=camera, got %q", q.Get("q"))
		}
		if q.Get("tag") != "electronics" {
			t.Errorf("expected tag=electronics, got %q", q.Get("tag"))
		}
		if q.Get("condition") != "excellent" {
			t.Errorf("expected condition=excellent, got %q", q.Get("condition"))
		}
		json.NewEncoder(w).Encode([]model.Item{})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.ListItems(store.ItemFilter{
		Query:     "camera",
		Tag:       "electronics",
		Condition: "excellent",
	})
	if err != nil {
		t.Fatalf("ListItems failed: %v", err)
	}
}

func TestClient_GetItem(t *testing.T) {
	id := uuid.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/items/"+id.String() {
			t.Errorf("expected path /api/items/%s, got %s", id, r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.Item{ID: id, Name: "Found"})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	item, err := c.GetItem(id)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}
	if item.Name != "Found" {
		t.Errorf("expected name 'Found', got %q", item.Name)
	}
}

func TestClient_CreateItem(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}

		var input model.ItemInput
		json.NewDecoder(r.Body).Decode(&input)
		if input.Name != "New Item" {
			t.Errorf("expected name 'New Item', got %q", input.Name)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(model.Item{ID: uuid.New(), Name: input.Name})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	item, err := c.CreateItem(model.ItemInput{Name: "New Item"})
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}
	if item.Name != "New Item" {
		t.Errorf("expected name 'New Item', got %q", item.Name)
	}
}

func TestClient_UpdateItem(t *testing.T) {
	id := uuid.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(model.Item{ID: id, Name: "Updated"})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	item, err := c.UpdateItem(id, model.ItemInput{Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdateItem failed: %v", err)
	}
	if item.Name != "Updated" {
		t.Errorf("expected name 'Updated', got %q", item.Name)
	}
}

func TestClient_DeleteItem(t *testing.T) {
	id := uuid.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	err := c.DeleteItem(id)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}
}

func TestClient_ListTags(t *testing.T) {
	tags := []model.Tag{{ID: 1, Name: "electronics"}, {ID: 2, Name: "tools"}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("expected /api/tags, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(tags)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	result, err := c.ListTags()
	if err != nil {
		t.Fatalf("ListTags failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result))
	}
}

func TestClient_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.GetItem(uuid.New())
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	if got := err.Error(); got != "HTTP 404: not found" {
		t.Errorf("expected 'HTTP 404: not found', got %q", got)
	}
}

func TestClient_UploadImage(t *testing.T) {
	// Create a temp file to upload
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.jpg")
	if err := os.WriteFile(testFile, []byte("fake jpeg data"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	itemID := uuid.New()
	imageID := uuid.New()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		f, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("failed to get form file: %v", err)
		}
		defer f.Close()

		if header.Filename != "test.jpg" {
			t.Errorf("expected filename test.jpg, got %q", header.Filename)
		}

		data, _ := io.ReadAll(f)
		if string(data) != "fake jpeg data" {
			t.Errorf("expected 'fake jpeg data', got %q", string(data))
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(model.Image{ID: imageID, ItemID: itemID, Filename: "test.jpg"})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	img, err := c.UploadImage(itemID, testFile)
	if err != nil {
		t.Fatalf("UploadImage failed: %v", err)
	}
	if img.Filename != "test.jpg" {
		t.Errorf("expected filename test.jpg, got %q", img.Filename)
	}
}

func TestClient_UploadImage_FileNotFound(t *testing.T) {
	c := New("http://localhost:0", "")
	_, err := c.UploadImage(uuid.New(), "/nonexistent/file.jpg")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestClient_DeleteImage(t *testing.T) {
	id := uuid.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/images/"+id.String() {
			t.Errorf("expected path /api/images/%s, got %s", id, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	if err := c.DeleteImage(id); err != nil {
		t.Fatalf("DeleteImage failed: %v", err)
	}
}
