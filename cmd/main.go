package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Jean1dev/communication-service/internal/infra"
	"github.com/joho/godotenv"
)

func init() {
	dir, _ := os.Getwd()
	log.Printf("Diretorio atual %v", dir)

	envPaths := []string{
		filepath.Join(dir, ".env"),
		filepath.Join(dir, "..", ".env"),
		filepath.Join(dir, "..", "..", ".env"),
	}

	var err error
	for _, envPath := range envPaths {
		err = godotenv.Load(envPath)
		if err == nil {
			log.Printf("Loaded .env from %s", envPath)
			return
		}
	}

	log.Printf("Warning: Could not load .env file from any of the attempted paths")
}

func main() {
	infra.ConfigAndStartHttpServer()
}
