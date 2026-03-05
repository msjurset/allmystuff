package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"allmystuff/internal/api"
	"allmystuff/internal/imgstore"
	"allmystuff/internal/secret"
	"allmystuff/internal/store"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	dbURL := envOr("ALLMYSTUFF_DB_URL", "postgres://localhost:5432/allmystuff?sslmode=disable")
	listen := envOr("ALLMYSTUFF_LISTEN", ":8080")
	imageDir := envOr("ALLMYSTUFF_IMAGE_DIR", defaultImageDir())

	// Run migrations
	if err := runMigrations(dbURL); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	s := store.NewPostgresStore(pool)
	imgs := imgstore.New(imageDir)

	apiKey, err := secret.Resolve(os.Getenv("ALLMYSTUFF_API_KEY"))
	if err != nil {
		log.Fatalf("resolving API key: %v", err)
	}
	if apiKey == "" {
		log.Println("WARNING: ALLMYSTUFF_API_KEY not set, API auth disabled")
	}

	router := api.NewRouter(s, imgs, apiKey)

	log.Printf("listening on %s", listen)
	if err := http.ListenAndServe(listen, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func runMigrations(dbURL string) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("migration source: %w", err)
	}

	// Ensure URL has proper scheme for migrate library
	migrateURL := dbURL
	if !strings.HasPrefix(migrateURL, "postgres://") && !strings.HasPrefix(migrateURL, "postgresql://") {
		migrateURL = "postgres://" + migrateURL
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, migrateURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func defaultImageDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".allmystuff", "images")
}
