package app

import (
	"go-scrfd-api/internal/config"
	"go-scrfd-api/internal/infrastructure/grpc"
	"go-scrfd-api/internal/infrastructure/rest"
	"go-scrfd-api/internal/recognition"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		config.NewConfig,
		grpc.NewGRPCServer,
		rest.NewRESTServer,
		recognition.NewUseCase,
		recognition.NewGRPCHandler,
		recognition.NewRESTHandler,
	),
	GrpcHandlers,
	fx.Invoke(grpc.RunGRPCServer, rest.RunRESTServer))
