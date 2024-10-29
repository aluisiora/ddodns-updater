package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    token := os.Getenv("DO_TOKEN")

    if (token == "") {
        log.Fatal("no token set")
    }

    record := os.Getenv("DNS_RECORD")

    if (record == "") {
        log.Fatal("no dns record set")
    }
}
