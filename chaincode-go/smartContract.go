package main

import (
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/handler"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/repository"
	"github.com/hyperledger/fabric-samples/asset-transfer-abac/chaincode-go/internal/service"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(repository.NewService),
		fx.Provide(service.NewService),		
		fx.Provide(handler.NewService),
		fx.Provide(internal.NewApp),
		fx.Invoke(func(*internal.App) {}),  // holds chaincode
	).Run()
}
