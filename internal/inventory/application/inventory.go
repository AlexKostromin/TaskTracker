package application

import (
	"context"

	"gitlab.com/godevs2/micro/internal/inventory/domain/model"
)

type InventoryStorage interface {
	Get(ctx context.Context, req *model.GetPartRequest) (*model.GetPartResponse, error)
	ListParts(ctx context.Context, req *model.ListPartsRequest) (*model.ListPartsResponse, error)
}

type InventoryService struct {
	inventoryRepo InventoryStorage
}

func NewInventoryService(inventoryRepo InventoryStorage) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) Get(ctx context.Context, req *model.GetPartRequest) (*model.GetPartResponse, error) {
	return s.inventoryRepo.Get(ctx, req)
}

func (s *InventoryService) ListParts(ctx context.Context, req *model.ListPartsRequest) (*model.ListPartsResponse, error) {
	return s.inventoryRepo.ListParts(ctx, req)
}
