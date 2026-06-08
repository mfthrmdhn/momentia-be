package endpoints

import (
	"momentia-be/internal/handler"
	"momentia-be/internal/middleware"
	"momentia-be/repository"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterPersonRoutes(r *gin.Engine, h handler.PersonHandler, hpd handler.PersonDateHandler, sessionRepo repository.UserSessionRepository) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET"), sessionRepo))
	{
		api.GET("/persons", h.GetPersons)
		api.GET("/persons/:id", h.GetPersonByID)
		api.POST("/persons", h.CreatePerson)
		api.PUT("/persons/:id", h.UpdatePerson)
		api.DELETE("/persons/:id", h.DeletePerson)

		api.GET("/persons/:id/date", hpd.GetAllPersonDates)
		api.GET("/persons/:id/date/:dateId", hpd.GetPersonDatesByID)
		api.POST("/persons/:id/date", hpd.CreatePersonDate)
		api.PUT("/persons/:id/date/:dateId", hpd.UpdatePersonDates)
		api.DELETE("/persons/:id/date/:dateId", hpd.DeletePersonDates)
	}
}
