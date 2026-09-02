package main

import (
	"log"

	"github.com/thegenggo/equipment-loan/api/internal/router"
)

func main() {
	r := router.Setup()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
