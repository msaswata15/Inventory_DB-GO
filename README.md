## Working (Project Progress & Details)

### Architecture
The project follows a clean architecture, separating concerns across multiple layers:
- **Domain:** Defines core entities (Category, Product, Order) and interfaces for repository and usecase contracts.
- **Repository:** Implements data access using GORM for PostgreSQL, handling CRUD operations and queries.
- **Usecase:** Contains business logic, including inventory management, stock validation, and order processing.
- **Service:** Provides asynchronous email notifications via SMTP, triggered on successful purchases.
- **Delivery:** Uses Gin for HTTP routing and request handling, mapping API endpoints to usecase methods.
- **Config:** Loads environment variables and application settings.
- **Database:** Manages GORM connection, schema migration, and database setup.

### Features Implemented
- **Category Management:**
  - List all categories (`GET /api/v1/collections`)
  - Create new category (`POST /api/v1/collections`)
- **Product Management:**
  - List all products (`GET /api/v1/items`)
  - Create new product (`POST /api/v1/items`)
- **Order Management:**
  - List all orders (`GET /api/v1/purchases`)
  - Create new order (`POST /api/v1/purchases`)
    - Reduces product stock atomically
    - Sends asynchronous email notification
- **Health Check:**
  - `GET /api/v1/pulse` returns server status

### Error Handling
- Returns 404 for empty lists (categories, products, orders)
- Returns 400 for insufficient stock on purchase
- Validates input payloads for required fields

### Email Notification
- Uses Go's `net/smtp` package
- Sends purchase confirmation emails asynchronously (goroutine)
- SMTP credentials loaded from environment

### Database
- PostgreSQL used for persistent storage
- GORM automigrates schema for categories, products, orders
- SQL migrations available in `migrations/`

### Environment & Configuration
- `.env.example` provided for environment setup
- Supports custom DB and SMTP settings

### Development Progress
- Core CRUD endpoints for categories, products, orders are functional
- Asynchronous email notifications implemented
- Error handling and validation in place
- Modular, testable code structure
- Further improvements planned: unit tests, API documentation, advanced filtering, admin endpoints
# Inventory Service Project (Go + Postgres + SMTP)

This repository is an ongoing project for an inventory management service, built with Go, PostgreSQL, and SMTP notifications. The design follows clean architecture principles and includes features for managing categories, products, and orders, with asynchronous email notifications. The API and structure are evolving as development progresses.

> **Note:** This is a work-in-progress project, not a course assignment. Features and structure may change as development continues.

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
