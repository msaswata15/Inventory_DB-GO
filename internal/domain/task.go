package domain

import (
	"context"
	"errors"
	"time"
)

// Category represents a logical grouping of products.
type Category struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" gorm:"type:text;unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product represents an item that can be purchased.
type Product struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name       string    `json:"name" gorm:"type:varchar(50);not null"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	CategoryID string    `json:"category_id" gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Order represents a purchase that reduces stock.
type Order struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskView is a read model combining orders, products, and categories.
type TaskView struct {
	OrderID      string    `json:"order_id"`
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	CategoryName string    `json:"category_name"`
	Quantity     int       `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
}

// InventoryRepository defines persistence operations.
type InventoryRepository interface {
	CreateCategory(ctx context.Context, c *Category) error
	ListCategories(ctx context.Context) ([]Category, error)
	CreateProduct(ctx context.Context, p *Product) error
	ListProducts(ctx context.Context) ([]Product, error)
	CreateOrder(ctx context.Context, o *Order) error
	ListOrders(ctx context.Context) ([]Order, error)
	DecrementStock(ctx context.Context, productID string, qty int) error
	GetTasks(ctx context.Context) ([]TaskView, error)
}

// InventoryUsecase defines business rules.
type InventoryUsecase interface {
	CreateCategory(ctx context.Context, name string) (*Category, error)
	ListCategories(ctx context.Context) ([]Category, error)
	CreateProduct(ctx context.Context, p *Product) (*Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
	CreateOrder(ctx context.Context, productID string, qty int) (*Order, error)
	ListOrders(ctx context.Context) ([]Order, error)
	GetTasks(ctx context.Context) ([]TaskView, error)
}

// EmailNotifier abstracts async notifications.
type EmailNotifier interface {
	SendAsync(subject, body string)
}

var (
	ErrNotFound          = errors.New("record not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)
