package main

import (
	"fmt"
	"github.com/syncrepair/backend/config"
)

func main() {
	cfg := config.Init()

	fmt.Println(cfg.AppName)
}
