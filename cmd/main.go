package main

import (
	"go-scrfd-api/internal/app"
	"go.uber.org/fx"
)

func main() {
	appRunner := fx.New(
		app.App,
	)
	appRunner.Run()
}
