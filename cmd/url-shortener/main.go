package main

import (
	"fmt"
	"log"
	"url-shortener/internal/config"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.MustLoad()

	fmt.Print(cfg)

	// TODO: init logger: slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: run server:
}
