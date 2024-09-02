package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Access environment variables
	listenPort := os.Getenv("LISTEN_PORT")
	log.Println("Listening on port:", listenPort)
}
