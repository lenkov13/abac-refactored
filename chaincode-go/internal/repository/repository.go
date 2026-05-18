package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/domain"
	"go.uber.org/fx"
)

// Service handles all direct read/write operations against the ledger world state.
type Service struct{}

func NewService(lc fx.Lifecycle) *Service {
	return &Service{}
}

func stub(ctx context.Context) shim.ChaincodeStubInterface {
	return ctx.Value(domain.StubKey).(shim.ChaincodeStubInterface)
}

func (r *Service) Save(ctx context.Context, asset *domain.Asset) error {
	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}
	return stub(ctx).PutState(asset.ID, data)
}

func (r *Service) FindByID(ctx context.Context, id string) (*domain.Asset, error) {
	data, err := stub(ctx).GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if data == nil {
		return nil, fmt.Errorf("asset %s does not exist", id)
	}
	var asset domain.Asset
	if err = json.Unmarshal(data, &asset); err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *Service) Exists(ctx context.Context, id string) (bool, error) {
	data, err := stub(ctx).GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %w", err)
	}
	return data != nil, nil
}

func (r *Service) Delete(ctx context.Context, id string) error {
	return stub(ctx).DelState(id)
}

func (r *Service) FindAll(ctx context.Context) ([]*domain.Asset, error) {
	iter, err := stub(ctx).GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var assets []*domain.Asset
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var asset domain.Asset
		if err = json.Unmarshal(item.Value, &asset); err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	return assets, nil
}
