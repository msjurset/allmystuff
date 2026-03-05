# allmystuff

A personal inventory management system with a REST API server and CLI client. Track your belongings with tags, images, and detailed metadata.

## Features

- **Item management** — Full CRUD for inventory items with name, brand, model, serial number, purchase info, condition, and notes
- **Tagging** — Organize items with arbitrary tags (auto-created on first use)
- **Image support** — Attach multiple images per item with ordering, up to 20MB each
- **API key auth** — Optional Bearer token authentication for the API
- **CLI client** — Full-featured command-line interface with table and JSON output
- **Filtering** — Search items by text query, tag, or condition

## Install

```sh
# From source
make deploy

# Or build both binaries locally
make build
```

Requires Go 1.25+ and a PostgreSQL database.

## Server

Start the API server:

```sh
# Minimal (defaults to localhost:8080, local PostgreSQL)
go run ./cmd/server

# Full configuration
ALLMYSTUFF_DB_URL="postgres://user:pass@host:5432/allmystuff" \
ALLMYSTUFF_LISTEN=":8080" \
ALLMYSTUFF_IMAGE_DIR="/data/images" \
ALLMYSTUFF_API_KEY="your-secret-key" \
  go run ./cmd/server
```

### Server Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ALLMYSTUFF_DB_URL` | `postgres://localhost:5432/allmystuff?sslmode=disable` | PostgreSQL connection string |
| `ALLMYSTUFF_LISTEN` | `:8080` | Server listen address |
| `ALLMYSTUFF_IMAGE_DIR` | `~/.allmystuff/images` | Image storage directory |
| `ALLMYSTUFF_API_KEY` | *(unset — auth disabled)* | API key for Bearer token auth. Supports 1Password secret references (`op://vault/item/field`) |

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/items` | List items (query params: `q`, `tag`, `condition`) |
| `POST` | `/api/items` | Create item |
| `GET` | `/api/items/{id}` | Get item |
| `PUT` | `/api/items/{id}` | Update item (partial merge) |
| `DELETE` | `/api/items/{id}` | Delete item |
| `POST` | `/api/items/{id}/images` | Upload image (multipart `file` field) |
| `PUT` | `/api/items/{id}/images/order` | Reorder images |
| `GET` | `/api/images/{id}` | Serve image |
| `DELETE` | `/api/images/{id}` | Delete image |
| `GET` | `/api/tags` | List tags |

## CLI

The `stuff` command-line tool communicates with the API server.

### CLI Configuration

| Flag | Env Var | Default | Description |
|------|---------|---------|-------------|
| `--url` | `ALLMYSTUFF_URL` | `http://localhost:8080` | API base URL |
| `--api-key` | `ALLMYSTUFF_API_KEY` | *(none)* | API key. Supports 1Password secret references (`op://vault/item/field`) |
| `--json` | — | `false` | Output raw JSON |

### Commands

#### Items

```sh
# List all items
stuff item list

# Search items
stuff item list -q "camera" --tag electronics --condition excellent

# Add an item
stuff item add --name "Sony A7IV" --brand Sony --model A7IV \
  --condition excellent --purchase-price 2499.99 \
  --purchase-date 2024-01-15 --tag camera,electronics

# Show item details
stuff item show <id>

# Edit an item (only specified fields are changed)
stuff item edit <id> --condition good --notes "minor scratch"

# Delete an item
stuff item delete <id>
stuff item delete <id> --yes  # skip confirmation
```

#### Images

```sh
# Upload an image
stuff image add <item-id> photo.jpg

# Delete an image
stuff image delete <image-id>
```

#### Tags

```sh
# List all tags
stuff tag list
```

### Item Flags

| Flag | Description |
|------|-------------|
| `--name` | Item name (required for `add`) |
| `--description` | Description |
| `--brand` | Brand |
| `--model` | Model |
| `--serial` | Serial number |
| `--purchase-date` | Purchase date (`YYYY-MM-DD`) |
| `--purchase-price` | Purchase price |
| `--estimated-value` | Estimated current value |
| `--condition` | Condition |
| `--notes` | Notes |
| `--tag` | Tags (comma-separated, repeatable) |

## Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE allmystuff;
```

Migrations run automatically on server startup.

## Build

```sh
make build       # Build both binaries
make test        # Run tests
make release     # Cross-compile for all platforms
make clean       # Remove build artifacts
```

## License

MIT — see [LICENSE](LICENSE).
