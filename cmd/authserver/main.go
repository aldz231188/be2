package main

import (
	"be2/internal/di"
	"context"
	"go.uber.org/fx"
	"log"
	"time"
)

func main() {
	startCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	app := fx.New(di.AuthApp)
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
