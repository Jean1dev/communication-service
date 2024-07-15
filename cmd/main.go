package main

import (
	"log"
	"os"

	"github.com/Jean1dev/communication-service/internal/infra"
	"github.com/joho/godotenv"
)

func init() {
	dir, _ := os.Getwd()
	log.Printf("Diretorio atual %v", dir)

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file %s", err.Error())
	}
}

func main() {
	infra.ConfigAndStartHttpServer()
}
