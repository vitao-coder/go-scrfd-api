package rest

import (
	"context"
	"go.uber.org/fx"
	"net/http"
)

type Handler interface {
	http.Handler
	Pattern() string
}

type ServerParams struct {
	fx.In
	Handlers []Handler `group:"rest_handlers"`
}

func NewRESTServer(p ServerParams) http.Handler {
	mux := http.NewServeMux()
	for _, h := range p.Handlers {
		mux.Handle("/", h)
	}
	return mux
}

func RunRESTServer(lc fx.Lifecycle, server http.Handler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go http.ListenAndServe(":8080", server)
			return nil
		},
	})
}
