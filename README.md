# Week 5 – Inventory Service (Go + Postgres + SMTP)

Refactored to mirror the inventory-management design (categories, products, orders) with clean architecture, GORM, and asynchronous SMTP notifications. Endpoint names are slightly adjusted from the reference to avoid exact matches.

## Folder layout
```
cmd/server          # App entrypoint
internal/config     # Env config loader
internal/database   # GORM connection + automigrate
internal/domain     # Entities + interfaces (contracts)
internal/repository # Postgres implementation (GORM)
internal/usecase    # Business logic
internal/service    # Email notifier
internal/delivery   # Gin handlers
migrations          # SQL for schema
```

## Prerequisites
- Go 1.21+
- PostgreSQL running locally
- Optional: SMTP creds (Gmail app password or similar)

## Setup
1) Copy `.env.example` to `.env` and fill DB + SMTP values.
2) Apply SQL (or rely on GORM automigrate):
   ```bash
   psql -U postgres -d inventory_db -f migrations/002_inventory.sql
   ```

## Run
```bash
cd "c:/Users/msasw/AMI/Week 5"
go run ./cmd/server
```
Server binds to `SERVER_PORT` (default 8080).

## API (renamed slightly)
- Health: `GET /api/v1/pulse`
- Categories ("collections"):
  - `GET /api/v1/collections`
  - `POST /api/v1/collections` `{ "name": "Laptops" }`
- Products ("items"):
  - `GET /api/v1/items`
  - `POST /api/v1/items` `{ "name": "XPS 13", "quantity": 5, "price_cents": 119900, "category_id": "<uuid>" }`
- Orders ("purchases") reduce stock atomically and send async email:
  - `GET /api/v1/purchases`
  - `POST /api/v1/purchases` `{ "product_id": "<uuid>", "quantity": 1 }`

Behavior notes:
- 404 when no records are present for list endpoints.
- 400 when stock is insufficient for a purchase.
- SMTP notifications fire in a goroutine; adjust `SMTP_*` vars to enable.

## Tech choices
- GORM for Postgres, auto-migrate categories/products/orders.
- Gin for HTTP delivery.
- SMTP via `net/smtp` with async goroutine.
- Clean layering: domain → repository → usecase → delivery.
