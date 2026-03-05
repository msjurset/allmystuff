package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"allmystuff/internal/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) ListItems(ctx context.Context, filter ItemFilter) ([]model.Item, error) {
	query := `SELECT DISTINCT i.id, i.name, i.description, i.brand, i.model, i.serial_number,
		i.purchase_date, i.purchase_price, i.estimated_value, i.condition, i.notes,
		i.created_at, i.updated_at
		FROM items i`

	var joins []string
	var conditions []string
	var args []any
	argIdx := 1

	if filter.Tag != "" {
		joins = append(joins, "JOIN item_tags it ON i.id = it.item_id JOIN tags t ON it.tag_id = t.id")
		conditions = append(conditions, fmt.Sprintf("t.name = $%d", argIdx))
		args = append(args, filter.Tag)
		argIdx++
	}

	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(i.name ILIKE $%d OR i.description ILIKE $%d OR i.brand ILIKE $%d OR i.model ILIKE $%d)", argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.Condition != "" {
		conditions = append(conditions, fmt.Sprintf("i.condition = $%d", argIdx))
		args = append(args, filter.Condition)
		argIdx++
	}

	if len(joins) > 0 {
		query += " " + strings.Join(joins, " ")
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY i.name"

	var items []model.Item
	if err := pgxscan.Select(ctx, s.pool, &items, query, args...); err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	if items == nil {
		items = []model.Item{}
	}

	for i := range items {
		if err := s.loadItemRelations(ctx, &items[i]); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *PostgresStore) GetItem(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	var item model.Item
	err := pgxscan.Get(ctx, s.pool, &item,
		`SELECT id, name, description, brand, model, serial_number,
			purchase_date, purchase_price, estimated_value, condition, notes,
			created_at, updated_at
		FROM items WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}
	if err := s.loadItemRelations(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *PostgresStore) CreateItem(ctx context.Context, input model.ItemInput) (*model.Item, error) {
	id := uuid.New()
	now := time.Now()

	var purchaseDate *time.Time
	if input.PurchaseDate != nil && *input.PurchaseDate != "" {
		t, err := time.Parse("2006-01-02", *input.PurchaseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid purchase_date: %w", err)
		}
		purchaseDate = &t
	}

	_, err := s.pool.Exec(ctx,
		`INSERT INTO items (id, name, description, brand, model, serial_number,
			purchase_date, purchase_price, estimated_value, condition, notes,
			created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		id, input.Name, input.Description, input.Brand, input.Model, input.SerialNumber,
		purchaseDate, input.PurchasePrice, input.EstimatedValue, input.Condition, input.Notes,
		now, now)
	if err != nil {
		return nil, fmt.Errorf("create item: %w", err)
	}

	if err := s.syncTags(ctx, id, input.Tags); err != nil {
		return nil, err
	}

	return s.GetItem(ctx, id)
}

func (s *PostgresStore) UpdateItem(ctx context.Context, id uuid.UUID, input model.ItemInput) (*model.Item, error) {
	var purchaseDate *time.Time
	if input.PurchaseDate != nil && *input.PurchaseDate != "" {
		t, err := time.Parse("2006-01-02", *input.PurchaseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid purchase_date: %w", err)
		}
		purchaseDate = &t
	}

	_, err := s.pool.Exec(ctx,
		`UPDATE items SET name=$1, description=$2, brand=$3, model=$4, serial_number=$5,
			purchase_date=$6, purchase_price=$7, estimated_value=$8, condition=$9, notes=$10,
			updated_at=NOW()
		WHERE id=$11`,
		input.Name, input.Description, input.Brand, input.Model, input.SerialNumber,
		purchaseDate, input.PurchasePrice, input.EstimatedValue, input.Condition, input.Notes,
		id)
	if err != nil {
		return nil, fmt.Errorf("update item: %w", err)
	}

	if err := s.syncTags(ctx, id, input.Tags); err != nil {
		return nil, err
	}

	return s.GetItem(ctx, id)
}

func (s *PostgresStore) DeleteItem(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}

func (s *PostgresStore) ListTags(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	if err := pgxscan.Select(ctx, s.pool, &tags, `SELECT id, name FROM tags ORDER BY name`); err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	if tags == nil {
		tags = []model.Tag{}
	}
	return tags, nil
}

func (s *PostgresStore) CreateImage(ctx context.Context, img model.Image) (*model.Image, error) {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO images (id, item_id, filename, filepath, sort_order, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		img.ID, img.ItemID, img.Filename, img.Filepath, img.SortOrder, img.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create image: %w", err)
	}
	return s.GetImage(ctx, img.ID)
}

func (s *PostgresStore) GetImage(ctx context.Context, id uuid.UUID) (*model.Image, error) {
	var img model.Image
	err := pgxscan.Get(ctx, s.pool, &img,
		`SELECT id, item_id, filename, filepath, sort_order, created_at FROM images WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}
	img.URL = fmt.Sprintf("/api/images/%s", img.ID)
	return &img, nil
}

func (s *PostgresStore) DeleteImage(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM images WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	return nil
}

func (s *PostgresStore) ListImages(ctx context.Context, itemID uuid.UUID) ([]model.Image, error) {
	var images []model.Image
	if err := pgxscan.Select(ctx, s.pool, &images,
		`SELECT id, item_id, filename, filepath, sort_order, created_at FROM images WHERE item_id = $1 ORDER BY sort_order`, itemID); err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	if images == nil {
		images = []model.Image{}
	}
	for i := range images {
		images[i].URL = fmt.Sprintf("/api/images/%s", images[i].ID)
	}
	return images, nil
}

func (s *PostgresStore) ReorderImages(ctx context.Context, itemID uuid.UUID, imageIDs []uuid.UUID) error {
	for i, imgID := range imageIDs {
		_, err := s.pool.Exec(ctx,
			`UPDATE images SET sort_order = $1 WHERE id = $2 AND item_id = $3`,
			i, imgID, itemID)
		if err != nil {
			return fmt.Errorf("reorder images: %w", err)
		}
	}
	return nil
}

// helpers

func (s *PostgresStore) loadItemRelations(ctx context.Context, item *model.Item) error {
	var tags []model.Tag
	if err := pgxscan.Select(ctx, s.pool, &tags,
		`SELECT t.id, t.name FROM tags t JOIN item_tags it ON t.id = it.tag_id WHERE it.item_id = $1 ORDER BY t.name`, item.ID); err != nil {
		return fmt.Errorf("load tags: %w", err)
	}
	if tags == nil {
		tags = []model.Tag{}
	}
	item.Tags = tags

	var images []model.Image
	if err := pgxscan.Select(ctx, s.pool, &images,
		`SELECT id, item_id, filename, filepath, sort_order, created_at FROM images WHERE item_id = $1 ORDER BY sort_order`, item.ID); err != nil {
		return fmt.Errorf("load images: %w", err)
	}
	if images == nil {
		images = []model.Image{}
	}
	for i := range images {
		images[i].URL = fmt.Sprintf("/api/images/%s", images[i].ID)
	}
	item.Images = images

	return nil
}

func (s *PostgresStore) getOrCreateTag(ctx context.Context, name string) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx, `SELECT id FROM tags WHERE name = $1`, name).Scan(&id)
	if err == nil {
		return id, nil
	}
	err = s.pool.QueryRow(ctx, `INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id`, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("get or create tag: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) syncTags(ctx context.Context, itemID uuid.UUID, tagNames []string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM item_tags WHERE item_id = $1`, itemID)
	if err != nil {
		return fmt.Errorf("clear item tags: %w", err)
	}
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		tagID, err := s.getOrCreateTag(ctx, name)
		if err != nil {
			return err
		}
		_, err = s.pool.Exec(ctx, `INSERT INTO item_tags (item_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, itemID, tagID)
		if err != nil {
			return fmt.Errorf("link tag: %w", err)
		}
	}
	return nil
}

