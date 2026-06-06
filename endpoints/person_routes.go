package endpoints

import (
	"momentia-be/internal/handler"
	"momentia-be/internal/middleware"
	"momentia-be/repository"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterPersonRoutes(r *gin.Engine, h handler.PersonHandler, sessionRepo repository.UserSessionRepository) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET"), sessionRepo))
	{
		api.GET("/persons", h.GetPersons)
		api.GET("/persons/:id", h.GetPersonByID)
		api.POST("/persons", h.CreatePerson)
		api.DELETE("/persons/:id", h.DeletePerson)
	}
}
