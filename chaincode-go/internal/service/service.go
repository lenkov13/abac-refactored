package service

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/domain"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/repository"
	"go.uber.org/fx"
)

// Service contains all business rules for asset management.
type Service struct {
	repo *repository.Service
}

func NewService(lc fx.Lifecycle, repo *repository.Service) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAsset(ctx context.Context, callerID string, hasCreatorAttr bool, mspID string, id, color string, size, appraisedValue int) error {
	// abac.creator is required to create asset
	if !hasCreatorAttr {
		return fmt.Errorf("submitting client not authorized to create asset, does not have abac.creator role")
	}
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}
	asset := &domain.Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          callerID,
		AppraisedValue: appraisedValue,
	}
	return s.repo.Save(ctx, asset)
}

func (s *Service) UpdateAsset(ctx context.Context, callerID, id, color string, size, value int) error {
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if callerID != asset.Owner {
		return fmt.Errorf("submitting client not authorized to update asset, does not own asset")
	}
	asset.Color = color
	asset.Size = size
	asset.AppraisedValue = value
	return s.repo.Save(ctx, asset)
}

func (s *Service) DeleteAsset(ctx context.Context, callerID, id string) error {
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if callerID != asset.Owner {
		return fmt.Errorf("submitting client not authorized to delete asset, does not own asset")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) TransferAsset(ctx context.Context, callerID, id, newOwner string) error {
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if callerID != asset.Owner {
		return fmt.Errorf("submitting client not authorized to transfer asset, does not own asset")
	}
	asset.Owner = newOwner
	return s.repo.Save(ctx, asset)
}

func (s *Service) ReadAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) GetAllAssets(ctx context.Context) ([]*domain.Asset, error) {
	return s.repo.FindAll(ctx)
}
