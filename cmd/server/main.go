package main

import (
	"log"
	"os"

	"github.com/aomori446/zuon/backend/api"
)

func main() {
	apiKey := os.Getenv("UNSPLASH_ACCESS_KEY")

	server, err := api.NewServer(apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	log.Println("Server starting on http://localhost:8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}