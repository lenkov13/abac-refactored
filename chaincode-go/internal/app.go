package internal

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/handler"
	"go.uber.org/fx"
)


// contractapi.NewChaincode(h) - wraps handler to Fabric compatible object
// Fanric will be calling handlers' methods (CreateAsset etc.) when
// peer sends the transaction  


// App owns the chaincode lifecycle managed by fx.
type App struct {
	chaincode *contractapi.ContractChaincode
}

func NewApp(lc fx.Lifecycle, h *handler.Service) *App {
	app := &App{}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			app.chaincode, err = contractapi.NewChaincode(h)
			if err != nil {
				return fmt.Errorf("error creating chaincode: %w", err)
			}
			// chaincode.Start() opens gRPC connection with peer
			// and begins listening for incoming transaction
			if err = app.chaincode.Start(); err != nil {
				return fmt.Errorf("error starting chaincode: %w", err)
			}
			return nil
		},
	})
	return app
}
