package store

import (
	"context"

	"allmystuff/internal/model"

	"github.com/google/uuid"
)

type ItemFilter struct {
	Query     string
	Tag       string
	Condition string
}

type Store interface {
	// Items
	ListItems(ctx context.Context, filter ItemFilter) ([]model.Item, error)
	GetItem(ctx context.Context, id uuid.UUID) (*model.Item, error)
	CreateItem(ctx context.Context, input model.ItemInput) (*model.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, input model.ItemInput) (*model.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error

	// Tags
	ListTags(ctx context.Context) ([]model.Tag, error)

	// Images
	CreateImage(ctx context.Context, img model.Image) (*model.Image, error)
	GetImage(ctx context.Context, id uuid.UUID) (*model.Image, error)
	DeleteImage(ctx context.Context, id uuid.UUID) error
	ListImages(ctx context.Context, itemID uuid.UUID) ([]model.Image, error)
	ReorderImages(ctx context.Context, itemID uuid.UUID, imageIDs []uuid.UUID) error
}
