package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"allmystuff/internal/model"
	"allmystuff/internal/store"

	"github.com/google/uuid"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTP:    &http.Client{},
	}
}

func (c *Client) do(method, path string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, bytes.TrimSpace(msg))
	}
	return resp, nil
}

func (c *Client) getJSON(path string, v any) error {
	resp, err := c.do("GET", path, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) getRaw(path string) ([]byte, error) {
	resp, err := c.do("GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) postJSON(path string, input, result any) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	resp, err := c.do("POST", path, bytes.NewReader(body), "application/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) putJSON(path string, input, result any) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	resp, err := c.do("PUT", path, bytes.NewReader(body), "application/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) delete(path string) error {
	resp, err := c.do("DELETE", path, nil, "")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// Items

func (c *Client) ListItems(filter store.ItemFilter) ([]model.Item, error) {
	params := url.Values{}
	if filter.Query != "" {
		params.Set("q", filter.Query)
	}
	if filter.Tag != "" {
		params.Set("tag", filter.Tag)
	}
	if filter.Condition != "" {
		params.Set("condition", filter.Condition)
	}

	path := "/api/items"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var items []model.Item
	if err := c.getJSON(path, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) ListItemsRaw(filter store.ItemFilter) ([]byte, error) {
	params := url.Values{}
	if filter.Query != "" {
		params.Set("q", filter.Query)
	}
	if filter.Tag != "" {
		params.Set("tag", filter.Tag)
	}
	if filter.Condition != "" {
		params.Set("condition", filter.Condition)
	}

	path := "/api/items"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	return c.getRaw(path)
}

func (c *Client) GetItem(id uuid.UUID) (*model.Item, error) {
	var item model.Item
	if err := c.getJSON(fmt.Sprintf("/api/items/%s", id), &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *Client) GetItemRaw(id uuid.UUID) ([]byte, error) {
	return c.getRaw(fmt.Sprintf("/api/items/%s", id))
}

func (c *Client) CreateItem(input model.ItemInput) (*model.Item, error) {
	var item model.Item
	if err := c.postJSON("/api/items", input, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *Client) CreateItemRaw(input model.ItemInput) ([]byte, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.do("POST", "/api/items", bytes.NewReader(body), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) UpdateItem(id uuid.UUID, input model.ItemInput) (*model.Item, error) {
	var item model.Item
	if err := c.putJSON(fmt.Sprintf("/api/items/%s", id), input, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *Client) UpdateItemRaw(id uuid.UUID, input model.ItemInput) ([]byte, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.do("PUT", fmt.Sprintf("/api/items/%s", id), bytes.NewReader(body), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) DeleteItem(id uuid.UUID) error {
	return c.delete(fmt.Sprintf("/api/items/%s", id))
}

// Tags

func (c *Client) ListTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := c.getJSON("/api/tags", &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (c *Client) ListTagsRaw() ([]byte, error) {
	return c.getRaw("/api/tags")
}

// Images

func (c *Client) UploadImage(itemID uuid.UUID, filePath string) (*model.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}
	w.Close()

	resp, err := c.do("POST", fmt.Sprintf("/api/items/%s/images", itemID), &buf, w.FormDataContentType())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var img model.Image
	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		return nil, err
	}
	return &img, nil
}

func (c *Client) DeleteImage(id uuid.UUID) error {
	return c.delete(fmt.Sprintf("/api/images/%s", id))
}
