package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thegenggo/equipment-loan/api/internal/handler"
)

func Setup(db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	healthHandler := handler.NewHealthHandler(db)
	r.GET("/health", healthHandler.Check)

	return r
}
