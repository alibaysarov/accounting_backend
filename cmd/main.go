package main

import (
	"acc_backend/internal/app"
	"log"
)

func main() {
	app := app.NewApp()

	if err := app.Run(); err != nil {
		log.Fatalf("An error: %v", err)
	}
}
