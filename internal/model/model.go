package model

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Description    string     `json:"description" db:"description"`
	Brand          string     `json:"brand" db:"brand"`
	Model          string     `json:"model" db:"model"`
	SerialNumber   string     `json:"serial_number" db:"serial_number"`
	PurchaseDate   *time.Time `json:"purchase_date,omitempty" db:"purchase_date"`
	PurchasePrice  *float64   `json:"purchase_price,omitempty" db:"purchase_price"`
	EstimatedValue *float64   `json:"estimated_value,omitempty" db:"estimated_value"`
	Condition      string     `json:"condition" db:"condition"`
	Notes          string     `json:"notes" db:"notes"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	Tags   []Tag   `json:"tags,omitempty"`
	Images []Image `json:"images,omitempty"`
}

type ItemInput struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Brand          string   `json:"brand"`
	Model          string   `json:"model"`
	SerialNumber   string   `json:"serial_number"`
	PurchaseDate   *string  `json:"purchase_date,omitempty"`
	PurchasePrice  *float64 `json:"purchase_price,omitempty"`
	EstimatedValue *float64 `json:"estimated_value,omitempty"`
	Condition      string   `json:"condition"`
	Notes          string   `json:"notes"`
	Tags []string `json:"tags,omitempty"`
}

type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Image struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ItemID    uuid.UUID `json:"item_id" db:"item_id"`
	Filename  string    `json:"filename" db:"filename"`
	Filepath  string    `json:"-" db:"filepath"`
	SortOrder int       `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	URL       string    `json:"url,omitempty" db:"-"`
}
