package endpoints

import (
	"momentia-be/internal/handler"
	"momentia-be/internal/middleware"
	"momentia-be/repository"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h handler.UserHandler, sessionRepo repository.UserSessionRepository) {
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET"), sessionRepo))
	{
		protected.GET("/profile", h.GetUserByID)
		protected.POST("/logout", h.Logout)
	}
}
