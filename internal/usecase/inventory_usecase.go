package usecase

import (
	"context"
	"fmt"
	"time"

	"week5/internal/domain"
)

// InventoryUsecase implements business rules around categories, products, and orders.
type InventoryUsecase struct {
	repo     domain.InventoryRepository
	notifier domain.EmailNotifier
}

func NewInventoryUsecase(repo domain.InventoryRepository, notifier domain.EmailNotifier) *InventoryUsecase {
	return &InventoryUsecase{repo: repo, notifier: notifier}
}

func (u *InventoryUsecase) CreateCategory(ctx context.Context, name string) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cat := &domain.Category{Name: name, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if err := u.repo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}

	body := fmt.Sprintf(
		"New category created:\n\nID: %s\nName: %s\nCreatedAt: %s\n",
		cat.ID, cat.Name, cat.CreatedAt.Format(time.RFC3339),
	)
	u.notifier.SendAsync("[Inventory] Category created", body)
	return cat, nil
}

func (u *InventoryUsecase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return u.repo.ListCategories(ctx)
}

func (u *InventoryUsecase) CreateProduct(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = p.CreatedAt

	if err := u.repo.CreateProduct(ctx, p); err != nil {
		return nil, err
	}
	body := fmt.Sprintf(
		"New product created:\n\nID: %s\nName: %s\nPrice: %.2f\nQuantity: %d\nCategoryID: %s\nCreatedAt: %s\n",
		p.ID, p.Name, p.Price, p.Quantity, p.CategoryID, p.CreatedAt.Format(time.RFC3339),
	)
	u.notifier.SendAsync("[Inventory] Product created", body)
	return p, nil
}

func (u *InventoryUsecase) ListProducts(ctx context.Context) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return u.repo.ListProducts(ctx)
}

func (u *InventoryUsecase) CreateOrder(ctx context.Context, productID string, qty int) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	order := &domain.Order{ProductID: productID, Quantity: qty, CreatedAt: time.Now().UTC()}
	if err := u.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}
	body := fmt.Sprintf(
		"Congratulations on your purchase!\n\n"+
			"We have successfully recorded your order with the following details:\n"+
			"- Order ID: %s\n"+
			"- Product ID: %s\n"+
			"- Quantity: %d\n"+
			"- Placed At: %s\n\n"+
			"Thank you for shopping with us!",
		order.ID, order.ProductID, order.Quantity, order.CreatedAt.Format(time.RFC3339),
	)
	u.notifier.SendAsync("Congratulations on your purchase", body)
	return order, nil
}

func (u *InventoryUsecase) ListOrders(ctx context.Context) ([]domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return u.repo.ListOrders(ctx)
}

func (u *InventoryUsecase) GetTasks(ctx context.Context) ([]domain.TaskView, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return u.repo.GetTasks(ctx)
}
