package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"week5/internal/domain"
)

type InventoryHandler struct {
	usecase domain.InventoryUsecase
}

func NewInventoryHandler(usecase domain.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{usecase: usecase}
}

func (h *InventoryHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	api.GET("/pulse", h.health)
	api.GET("/tasks", h.getTasks)

	api.GET("/collections", h.listCategories)
	api.POST("/collections", h.createCategory)

	api.GET("/items", h.listProducts)
	api.POST("/items", h.createProduct)

	api.GET("/purchases", h.listOrders)
	api.POST("/purchases", h.createOrder)
}

func (h *InventoryHandler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func (h *InventoryHandler) listCategories(c *gin.Context) {
	categories, err := h.usecase.ListCategories(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *InventoryHandler) createCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.usecase.CreateCategory(c.Request.Context(), req.Name)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": cat})
}

func (h *InventoryHandler) listProducts(c *gin.Context) {
	products, err := h.usecase.ListProducts(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *InventoryHandler) createProduct(c *gin.Context) {
	var req struct {
		Name       string  `json:"name" binding:"required"`
		Quantity   int     `json:"quantity" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		CategoryID string  `json:"category_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &domain.Product{
		Name:       req.Name,
		Quantity:   req.Quantity,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}

	created, err := h.usecase.CreateProduct(c.Request.Context(), product)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func (h *InventoryHandler) listOrders(c *gin.Context) {
	orders, err := h.usecase.ListOrders(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *InventoryHandler) createOrder(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.usecase.CreateOrder(c.Request.Context(), req.ProductID, req.Quantity)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": order})
}

func (h *InventoryHandler) getTasks(c *gin.Context) {
	tasks, err := h.usecase.GetTasks(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (h *InventoryHandler) handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "no records"})
	case domain.ErrInsufficientStock:
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
