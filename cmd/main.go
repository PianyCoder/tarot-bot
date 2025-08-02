package main

import (
	"TarotBot/internal/app"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("приложение запустилось с ошибкой %w", err)
	}
}
