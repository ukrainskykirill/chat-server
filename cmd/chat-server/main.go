package main

import (
	"context"
	"github.com/ukrainskykirill/chat-server/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
