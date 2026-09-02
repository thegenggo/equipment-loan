package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/handler"
)

func Setup() *gin.Engine {
	r := gin.Default()

	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	return r
}
