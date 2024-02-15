package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/syncrepair/backend/internal/config"
)

func main() {
	// Configuration
	cfg := config.Load()

	fmt.Println(cfg)
}
