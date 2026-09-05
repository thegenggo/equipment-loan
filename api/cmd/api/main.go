package main

import (
	"log"

	"github.com/thegenggo/equipment-loan/api/internal/config"
	"github.com/thegenggo/equipment-loan/api/internal/database"
	"github.com/thegenggo/equipment-loan/api/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	r := router.Setup(db)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
