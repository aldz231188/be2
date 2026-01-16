package main

import (
	"go.uber.org/fx"

	"be2/services/clientsvc/di"
)

func main() {
	fx.New(di.App).Run()
}
