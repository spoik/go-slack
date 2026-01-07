package main

import (
	"context"
	"go-slack/database"
	"go-slack/httpserver"
	"log"
)

func main() {
	ctx := context.Background()
	db, err := database.Connect(ctx)

	if err != nil {
		log.Printf("Failed to connect to the database: %e\n", err)
		return
	}

	defer db.Close(ctx)
	httpserver.StartNew(ctx, db)
}
