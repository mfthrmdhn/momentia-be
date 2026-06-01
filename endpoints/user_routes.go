package endpoints

import (
	"momentia-be/internal/handler"
	"momentia-be/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h handler.UserHandler) {
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
	
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		protected.GET("/profile", h.GetUserByID)
		protected.POST("/logout", h.Logout)
	}
}