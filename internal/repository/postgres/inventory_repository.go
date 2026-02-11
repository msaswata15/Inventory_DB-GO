package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"week5/internal/domain"
)

// InventoryRepository is a GORM-backed implementation.
type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) CreateCategory(ctx context.Context, c *domain.Category) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *InventoryRepository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, domain.ErrNotFound
	}
	return categories, nil
}

func (r *InventoryRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *InventoryRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, domain.ErrNotFound
	}
	return products, nil
}

func (r *InventoryRepository) CreateOrder(ctx context.Context, o *domain.Order) error {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	o.CreatedAt = time.Now().UTC()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.decrementStockTx(ctx, tx, o.ProductID, o.Quantity); err != nil {
			return err
		}
		return tx.Create(o).Error
	})
}

func (r *InventoryRepository) ListOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, domain.ErrNotFound
	}
	return orders, nil
}

func (r *InventoryRepository) GetTasks(ctx context.Context) ([]domain.TaskView, error) {
	var tasks []domain.TaskView

	q := r.db.WithContext(ctx).
		Table("orders o").
		Select(`o.id as order_id, o.product_id, p.name as product_name, c.name as category_name, o.quantity, o.created_at`).
		Joins("JOIN products p ON p.id = o.product_id").
		Joins("JOIN categories c ON c.id = p.category_id").
		Order("o.created_at DESC")

	if err := q.Scan(&tasks).Error; err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, domain.ErrNotFound
	}
	return tasks, nil
}

// DecrementStock performs an atomic quantity decrement.
func (r *InventoryRepository) DecrementStock(ctx context.Context, productID string, qty int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.decrementStockTx(ctx, tx, productID, qty)
	})
}

func (r *InventoryRepository) decrementStockTx(ctx context.Context, tx *gorm.DB, productID string, qty int) error {
	if qty <= 0 {
		return errors.New("quantity must be positive")
	}

	res := tx.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? AND quantity >= ?", productID, qty).
		Update("quantity", gorm.Expr("quantity - ?", qty))

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrInsufficientStock
	}

	// Lock row to avoid race when reading updated quantity if needed later.
	return tx.WithContext(ctx).Clauses(clause.Locking{Strength: "SHARE"}).
		First(&domain.Product{}, "id = ?", productID).Error
}
