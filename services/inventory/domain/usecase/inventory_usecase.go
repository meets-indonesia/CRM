package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/inventory/domain/repository"
)

// Errors
var (
	ErrItemNotFound      = errors.New("item not found")
	ErrSupplierNotFound  = errors.New("supplier not found")
	ErrDuplicateSKU      = errors.New("item with this SKU already exists")
	ErrInsufficientStock = errors.New("insufficient stock")
)

// InventoryUsecase mendefinisikan operasi-operasi usecase untuk Inventory
type InventoryUsecase interface {
	// Item operations
	CreateItem(ctx context.Context, req entity.CreateItemRequest) (*entity.Item, error)
	GetItem(ctx context.Context, id uint) (*entity.Item, error)
	UpdateItem(ctx context.Context, id uint, req entity.UpdateItemRequest) (*entity.Item, error)
	DeleteItem(ctx context.Context, id uint) error
	ListItems(ctx context.Context, activeOnly bool, page, limit int) (*entity.ItemListResponse, error)
	ListItemsByCategory(ctx context.Context, category string, activeOnly bool, page, limit int) (*entity.ItemListResponse, error)
	SearchItems(ctx context.Context, query string, activeOnly bool, page, limit int) (*entity.ItemListResponse, error)
	GetLowStockItems(ctx context.Context) (*entity.LowStockResponse, error)

	// Stock operations
	IncreaseStock(ctx context.Context, itemID uint, req entity.StockUpdateRequest, performedBy uint) (*entity.StockTransaction, error)
	DecreaseStock(ctx context.Context, itemID uint, req entity.StockUpdateRequest, performedBy uint) (*entity.StockTransaction, error)
	GetStockTransaction(ctx context.Context, id uint) (*entity.StockTransaction, error)
	ListStockTransactionsByItem(ctx context.Context, itemID uint, page, limit int) (*entity.StockTransactionListResponse, error)
	ListAllStockTransactions(ctx context.Context, page, limit int) (*entity.StockTransactionListResponse, error)

	// Supplier operations
	CreateSupplier(ctx context.Context, req entity.CreateSupplierRequest) (*entity.Supplier, error)
	GetSupplier(ctx context.Context, id uint) (*entity.Supplier, error)
	UpdateSupplier(ctx context.Context, id uint, req entity.UpdateSupplierRequest) (*entity.Supplier, error)
	DeleteSupplier(ctx context.Context, id uint) error
	ListSuppliers(ctx context.Context, activeOnly bool, page, limit int) (*entity.SupplierListResponse, error)
	SearchSuppliers(ctx context.Context, query string, activeOnly bool, page, limit int) (*entity.SupplierListResponse, error)

	// Event Processing
	ProcessRewardClaimed(claimID uint, rewardID uint, quantity int) error
}

type inventoryUsecase struct {
	itemRepo       repository.ItemRepository
	stockRepo      repository.StockRepository
	supplierRepo   repository.SupplierRepository
	eventPublisher repository.EventPublisher
}

// NewInventoryUsecase membuat instance baru InventoryUsecase
func NewInventoryUsecase(
	itemRepo repository.ItemRepository,
	stockRepo repository.StockRepository,
	supplierRepo repository.SupplierRepository,
	eventPublisher repository.EventPublisher,
) InventoryUsecase {
	return &inventoryUsecase{
		itemRepo:       itemRepo,
		stockRepo:      stockRepo,
		supplierRepo:   supplierRepo,
		eventPublisher: eventPublisher,
	}
}

// CreateItem membuat item baru
func (u *inventoryUsecase) CreateItem(ctx context.Context, req entity.CreateItemRequest) (*entity.Item, error) {
	// Check if SKU already exists
	existingItem, err := u.itemRepo.FindBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}
	if existingItem != nil {
		return nil, ErrDuplicateSKU
	}

	// Check if supplier exists if provided
	if req.SupplierID != nil && *req.SupplierID > 0 {
		supplier, err := u.supplierRepo.FindByID(ctx, *req.SupplierID)
		if err != nil {
			return nil, err
		}
		if supplier == nil {
			return nil, ErrSupplierNotFound
		}
	}

	// Create item
	item := &entity.Item{
		Name:            req.Name,
		Description:     req.Description,
		SKU:             req.SKU,
		ImageURL:        req.ImageURL,
		Category:        req.Category,
		CurrentStock:    req.InitialStock,
		MinimumStock:    req.MinimumStock,
		SupplierID:      req.SupplierID,
		ReorderQuantity: req.ReorderQuantity,
		IsActive:        true,
	}

	if err := u.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	// If initial stock is greater than zero, create stock transaction
	if req.InitialStock > 0 {
		transaction := &entity.StockTransaction{
			ItemID:      item.ID,
			Type:        entity.TransactionTypeIncrease,
			Quantity:    req.InitialStock,
			PreviousQty: 0,
			NewQty:      req.InitialStock,
			Reason:      "Initial stock",
			PerformedBy: 1, // System user ID
			CreatedAt:   time.Now(),
			Item:        *item,
		}

		if err := u.stockRepo.CreateTransaction(ctx, transaction); err != nil {
			// Log error but don't fail
		}
	}

	// Check if stock is below minimum and alert if necessary
	if item.CurrentStock < item.MinimumStock {
		deficit := item.MinimumStock - item.CurrentStock
		if err := u.eventPublisher.PublishLowStockAlert(item, deficit); err != nil {
			// Log error but don't fail
		}
	}

	return item, nil
}

// GetItem mendapatkan item berdasarkan ID
func (u *inventoryUsecase) GetItem(ctx context.Context, id uint) (*entity.Item, error) {
	item, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	return item, nil
}

// UpdateItem memperbarui item
func (u *inventoryUsecase) UpdateItem(ctx context.Context, id uint, req entity.UpdateItemRequest) (*entity.Item, error) {
	item, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Check if new SKU already exists
	if req.SKU != "" && req.SKU != item.SKU {
		existingItem, err := u.itemRepo.FindBySKU(ctx, req.SKU)
		if err != nil {
			return nil, err
		}
		if existingItem != nil {
			return nil, ErrDuplicateSKU
		}
	}

	// Check if supplier exists if provided
	if req.SupplierID != nil && *req.SupplierID > 0 {
		supplier, err := u.supplierRepo.FindByID(ctx, *req.SupplierID)
		if err != nil {
			return nil, err
		}
		if supplier == nil {
			return nil, ErrSupplierNotFound
		}
	}

	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.SKU != "" {
		item.SKU = req.SKU
	}
	if req.ImageURL != "" {
		item.ImageURL = req.ImageURL
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.MinimumStock != nil {
		item.MinimumStock = *req.MinimumStock
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if req.SupplierID != nil {
		item.SupplierID = req.SupplierID
	}
	if req.ReorderQuantity != nil {
		item.ReorderQuantity = *req.ReorderQuantity
	}

	if err := u.itemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	// Check if stock is below minimum after update
	if item.CurrentStock < item.MinimumStock {
		deficit := item.MinimumStock - item.CurrentStock
		if err := u.eventPublisher.PublishLowStockAlert(item, deficit); err != nil {
			// Log error but don't fail
		}
	}

	return item, nil
}

// DeleteItem menghapus item
func (u *inventoryUsecase) DeleteItem(ctx context.Context, id uint) error {
	item, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrItemNotFound
	}

	return u.itemRepo.Delete(ctx, id)
}

// ListItems mendapatkan daftar item
func (u *inventoryUsecase) ListItems(ctx context.Context, activeOnly bool, page, limit int) (*entity.ItemListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	items, total, err := u.itemRepo.List(ctx, activeOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ItemListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// ListItemsByCategory mendapatkan daftar item berdasarkan kategori
func (u *inventoryUsecase) ListItemsByCategory(ctx context.Context, category string, activeOnly bool, page, limit int) (*entity.ItemListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	items, total, err := u.itemRepo.ListByCategory(ctx, category, activeOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ItemListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// SearchItems mencari item berdasarkan kata kunci
func (u *inventoryUsecase) SearchItems(ctx context.Context, query string, activeOnly bool, page, limit int) (*entity.ItemListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	items, total, err := u.itemRepo.Search(ctx, query, activeOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ItemListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// GetLowStockItems mendapatkan item dengan stok di bawah minimum
func (u *inventoryUsecase) GetLowStockItems(ctx context.Context) (*entity.LowStockResponse, error) {
	items, total, err := u.itemRepo.FindLowStockItems(ctx)
	if err != nil {
		return nil, err
	}

	lowStockItems := make([]entity.LowStockItem, len(items))
	for i, item := range items {
		deficit := item.MinimumStock - item.CurrentStock
		lowStockItems[i] = entity.LowStockItem{
			Item:          item,
			Deficit:       deficit,
			ReorderAmount: item.ReorderQuantity,
		}
	}

	return &entity.LowStockResponse{
		Items: lowStockItems,
		Total: total,
	}, nil
}

// IncreaseStock menambah stok item
func (u *inventoryUsecase) IncreaseStock(ctx context.Context, itemID uint, req entity.StockUpdateRequest, performedBy uint) (*entity.StockTransaction, error) {
	if req.Quantity <= 0 {
		req.Quantity = -req.Quantity // Make sure quantity is positive
	}

	item, err := u.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Create transaction first so we have a record even if update fails
	previousStock := item.CurrentStock
	newStock := previousStock + req.Quantity

	transaction := &entity.StockTransaction{
		ItemID:      itemID,
		Type:        entity.TransactionTypeIncrease,
		Quantity:    req.Quantity,
		PreviousQty: previousStock,
		NewQty:      newStock,
		Reason:      req.Reason,
		ReferenceID: req.ReferenceID,
		PerformedBy: performedBy,
		CreatedAt:   time.Now(),
		Item:        *item,
	}

	// Update item stock
	if err := u.itemRepo.UpdateStock(ctx, itemID, newStock); err != nil {
		return nil, err
	}

	// Save transaction
	if err := u.stockRepo.CreateTransaction(ctx, transaction); err != nil {
		// Revert stock update if transaction fails
		u.itemRepo.UpdateStock(ctx, itemID, previousStock)
		return nil, err
	}

	// Publish event
	if err := u.eventPublisher.PublishStockUpdated(transaction); err != nil {
		// Log error but don't fail
	}

	return transaction, nil
}

// DecreaseStock mengurangi stok item
func (u *inventoryUsecase) DecreaseStock(ctx context.Context, itemID uint, req entity.StockUpdateRequest, performedBy uint) (*entity.StockTransaction, error) {
	if req.Quantity <= 0 {
		req.Quantity = -req.Quantity // Make sure quantity is positive
	}

	item, err := u.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Check if stock is sufficient
	if item.CurrentStock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	// Create transaction first so we have a record even if update fails
	previousStock := item.CurrentStock
	newStock := previousStock - req.Quantity

	transaction := &entity.StockTransaction{
		ItemID:      itemID,
		Type:        entity.TransactionTypeDecrease,
		Quantity:    req.Quantity,
		PreviousQty: previousStock,
		NewQty:      newStock,
		Reason:      req.Reason,
		ReferenceID: req.ReferenceID,
		PerformedBy: performedBy,
		CreatedAt:   time.Now(),
		Item:        *item,
	}

	// Update item stock
	if err := u.itemRepo.UpdateStock(ctx, itemID, newStock); err != nil {
		return nil, err
	}

	// Save transaction
	if err := u.stockRepo.CreateTransaction(ctx, transaction); err != nil {
		// Revert stock update if transaction fails
		u.itemRepo.UpdateStock(ctx, itemID, previousStock)
		return nil, err
	}

	// Check if stock is below minimum and alert if necessary
	if newStock < item.MinimumStock {
		deficit := item.MinimumStock - newStock
		if err := u.eventPublisher.PublishLowStockAlert(item, deficit); err != nil {
			// Log error but don't fail
		}
	}

	// Publish event
	if err := u.eventPublisher.PublishStockUpdated(transaction); err != nil {
		// Log error but don't fail
	}

	return transaction, nil
}

// GetStockTransaction mendapatkan transaksi stok berdasarkan ID
func (u *inventoryUsecase) GetStockTransaction(ctx context.Context, id uint) (*entity.StockTransaction, error) {
	transaction, err := u.stockRepo.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errors.New("stock transaction not found")
	}

	return transaction, nil
}

// ListStockTransactionsByItem mendapatkan daftar transaksi stok berdasarkan item
func (u *inventoryUsecase) ListStockTransactionsByItem(ctx context.Context, itemID uint, page, limit int) (*entity.StockTransactionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, total, err := u.stockRepo.ListTransactionsByItemID(ctx, itemID, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.StockTransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
}

// ListAllStockTransactions mendapatkan semua transaksi stok
func (u *inventoryUsecase) ListAllStockTransactions(ctx context.Context, page, limit int) (*entity.StockTransactionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, total, err := u.stockRepo.ListAllTransactions(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.StockTransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
}

// CreateSupplier membuat supplier baru
func (u *inventoryUsecase) CreateSupplier(ctx context.Context, req entity.CreateSupplierRequest) (*entity.Supplier, error) {
	supplier := &entity.Supplier{
		Name:        req.Name,
		ContactName: req.ContactName,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		IsActive:    true,
	}

	if err := u.supplierRepo.Create(ctx, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// GetSupplier mendapatkan supplier berdasarkan ID
func (u *inventoryUsecase) GetSupplier(ctx context.Context, id uint) (*entity.Supplier, error) {
	supplier, err := u.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, ErrSupplierNotFound
	}

	return supplier, nil
}

// UpdateSupplier memperbarui supplier
func (u *inventoryUsecase) UpdateSupplier(ctx context.Context, id uint, req entity.UpdateSupplierRequest) (*entity.Supplier, error) {
	supplier, err := u.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, ErrSupplierNotFound
	}

	// Update fields
	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.ContactName != "" {
		supplier.ContactName = req.ContactName
	}
	if req.Email != "" {
		supplier.Email = req.Email
	}
	if req.Phone != "" {
		supplier.Phone = req.Phone
	}
	if req.Address != "" {
		supplier.Address = req.Address
	}
	if req.IsActive != nil {
		supplier.IsActive = *req.IsActive
	}

	if err := u.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// DeleteSupplier menghapus supplier
func (u *inventoryUsecase) DeleteSupplier(ctx context.Context, id uint) error {
	supplier, err := u.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if supplier == nil {
		return ErrSupplierNotFound
	}

	return u.supplierRepo.Delete(ctx, id)
}

// ListSuppliers mendapatkan daftar supplier
func (u *inventoryUsecase) ListSuppliers(ctx context.Context, activeOnly bool, page, limit int) (*entity.SupplierListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	suppliers, total, err := u.supplierRepo.List(ctx, activeOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.SupplierListResponse{
		Suppliers: suppliers,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// SearchSuppliers mencari supplier berdasarkan kata kunci
func (u *inventoryUsecase) SearchSuppliers(ctx context.Context, query string, activeOnly bool, page, limit int) (*entity.SupplierListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	suppliers, total, err := u.supplierRepo.Search(ctx, query, activeOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.SupplierListResponse{
		Suppliers: suppliers,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// ProcessRewardClaimed memproses event hadiah diklaim
func (u *inventoryUsecase) ProcessRewardClaimed(claimID uint, rewardID uint, quantity int) error {
	// Jika quantity tidak diberikan, default ke 1
	if quantity <= 0 {
		quantity = 1
	}

	// Cari item berdasarkan rewardID
	// Asumsikan rewardID sama dengan itemID, atau bisa ditambahkan mapping jika berbeda
	item, err := u.itemRepo.FindByID(context.Background(), rewardID)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrItemNotFound
	}

	// Buat permintaan pengurangan stok
	req := entity.StockUpdateRequest{
		Quantity:    quantity,
		Reason:      "Reward claimed",
		ReferenceID: "claim:" + string(claimID),
	}

	// Kurangi stok
	_, err = u.DecreaseStock(context.Background(), rewardID, req, 1) // System user ID
	if err != nil {
		return err
	}

	return nil
}
