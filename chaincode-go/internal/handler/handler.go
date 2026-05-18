package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/domain"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/service"
	"go.uber.org/fx"
)

// Service is the chaincode entry point; it translates TX context into service calls.
type Service struct {
	contractapi.Contract
	bl *service.Service
}

func NewService(lc fx.Lifecycle, bl *service.Service) *Service {
	return &Service{bl: bl}
}

// txContext extracts: 
// the stub
// decoded caller identity (base64 -> string)
// abac.creator attribute from the TX contex (need for ABAC contract)

func txContext(c contractapi.TransactionContextInterface) (ctx context.Context, callerID string, hasCreatorAttr bool, mspID string) {
	ctx = context.WithValue(context.Background(), domain.StubKey, c.GetStub())

	b64ID, err := c.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("failed to get client identity: %v", err)
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		log.Printf("failed to decode client identity: %v", err)
		return
	}
	callerID = string(decoded)

	hasCreatorAttr = c.GetClientIdentity().AssertAttributeValue("abac.creator", "true") == nil

	mspID, err = c.GetClientIdentity().GetMSPID()
	if err != nil {
		log.Printf("failed to get mspID: %v", err)
	}
	return
}

func (h *Service) CreateAsset(c contractapi.TransactionContextInterface, id, color string, size, appraisedValue int) error {
	ctx, callerID, hasCreatorAttr, mspID := txContext(c)
	return h.bl.CreateAsset(ctx, callerID, hasCreatorAttr, mspID, id, color, size, appraisedValue)
}

func (h *Service) UpdateAsset(c contractapi.TransactionContextInterface, id, newColor string, newSize, newValue int) error {
	ctx, callerID, _, _ := txContext(c)
	return h.bl.UpdateAsset(ctx, callerID, id, newColor, newSize, newValue)
}

func (h *Service) DeleteAsset(c contractapi.TransactionContextInterface, id string) error {
	ctx, callerID, _, _ := txContext(c)
	return h.bl.DeleteAsset(ctx, callerID, id)
}

func (h *Service) TransferAsset(c contractapi.TransactionContextInterface, id, newOwner string) error {
	ctx, callerID, _, _ := txContext(c)
	return h.bl.TransferAsset(ctx, callerID, id, newOwner)
}

func (h *Service) ReadAsset(c contractapi.TransactionContextInterface, id string) (*domain.Asset, error) {
	ctx, _, _, _ := txContext(c)
	return h.bl.ReadAsset(ctx, id)
}

func (h *Service) GetAllAssets(c contractapi.TransactionContextInterface) ([]*domain.Asset, error) {
	ctx, _, _, _ := txContext(c)
	return h.bl.GetAllAssets(ctx)
}

// GetSubmittingClientIdentity exposes the decoded caller identity as a chaincode function.
func (h *Service) GetSubmittingClientIdentity(c contractapi.TransactionContextInterface) (string, error) {
	_, callerID, _, _ := txContext(c)
	if callerID == "" {
		return "", fmt.Errorf("failed to get submitting client identity")
	}
	return callerID, nil
}
