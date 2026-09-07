package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thegenggo/equipment-loan/api/internal/config"
	"github.com/thegenggo/equipment-loan/api/internal/handler"
	"github.com/thegenggo/equipment-loan/api/internal/middleware"
	"github.com/thegenggo/equipment-loan/api/internal/model"
	"github.com/thegenggo/equipment-loan/api/internal/repository"
	"github.com/thegenggo/equipment-loan/api/internal/service"
	"github.com/thegenggo/equipment-loan/api/pkg/token"
)

func Setup(cfg *config.Config, db *sqlx.DB) *gin.Engine {
	tokens := token.NewManager(cfg.JWTSecret, cfg.JWTTTL)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, tokens)

	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authService)

	equipmentRepo := repository.NewEquipmentRepository(db)
	equipmentService := service.NewEquipmentService(equipmentRepo)
	equipmentHandler := handler.NewEquipmentHandler(equipmentService)

	r := gin.Default()

	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth(tokens))
		{
			protected.GET("/me", authHandler.Me)
		}

		equipments := api.Group("/equipments")
		equipments.Use(middleware.Auth(tokens))
		{
			equipments.GET("", equipmentHandler.List)
			equipments.GET("/:id", equipmentHandler.Get)

			admin := equipments.Group("")
			admin.Use(middleware.RequireRole(model.RoleAdmin))
			{
				admin.POST("", equipmentHandler.Create)
				admin.PUT("/:id", equipmentHandler.Update)
				admin.DELETE("/:id", equipmentHandler.Delete)
			}
		}
	}

	return r
}
