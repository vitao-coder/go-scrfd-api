package app

import (
	"go-scrfd-api/internal/infrastructure/grpc"
	"go-scrfd-api/internal/recognition"
	"go.uber.org/fx"
)

var GrpcHandlers = fx.Provide(NewRecognitionGrpcServiceHandler)

func NewRecognitionGrpcServiceHandler(uc recognition.UseCase) grpc.ServerOutput {
	return grpc.ServerOutput{
		Handler: recognition.NewGRPCHandler(uc),
	}
}
