package grpc

import (
	"context"
	"go-scrfd-api/pkg/grpc/recognition/pb"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
)

type ServerOutput struct {
	fx.Out
	Handler pb.RecognitionServer `group:"grpc_handlers"`
}

type ServerInput struct {
	fx.In
	Handlers []pb.RecognitionServer `group:"grpc_handlers"`
}

func NewGRPCServer(p ServerInput) *grpc.Server {
	s := grpc.NewServer()
	for _, h := range p.Handlers {
		pb.RegisterRecognitionServer(s, h)
	}
	return s
}

func RunGRPCServer(lc fx.Lifecycle, server *grpc.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":50051")
			if err != nil {
				return err
			}
			go server.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()
			return nil
		},
	})
}
