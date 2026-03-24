package routes

import (
	"stellarbill-backend/internal/config"
	"stellarbill-backend/internal/cors"
	"stellarbill-backend/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/handlers"
	"stellarbill-backend/internal/middleware"
	"stellarbill-backend/internal/repository"
	"stellarbill-backend/internal/service"
)

func Register(r *gin.Engine) {
	cfg := config.Load()
	corsProfile := cors.ProfileForEnv(cfg.Env, cfg.AllowedOrigins)

	r.Use(cors.Middleware(corsProfile))

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}

	subRepo := repository.NewMockSubscriptionRepo()
	planRepo := repository.NewMockPlanRepo()
	svc := service.NewSubscriptionService(subRepo, planRepo)

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
		api.GET("/subscriptions", handlers.ListSubscriptions)
		api.GET("/subscriptions/:id", middleware.AuthMiddleware(jwtSecret), handlers.NewGetSubscriptionHandler(svc))
		api.GET("/plans", handlers.ListPlans)
	}
}
