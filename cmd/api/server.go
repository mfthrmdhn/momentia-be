package main

import (
	"log"
	"momentia-be/internal/config"
	"momentia-be/internal/handler"
	"momentia-be/internal/middleware"
	"momentia-be/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func runMigrations(dsn string) {
	m, err := migrate.New("file://internal/db/migrations", dsn)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")
}

func main() {
	// Initialize the server
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Init db
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		panic("database connection failed: " + err.Error())
	}
	runMigrations(cfg.DB.DSN())

	// Wire up dependencies
	personRepo := repository.NewPersonRepository(db)
	personHandler := handler.NewPersonHandler(personRepo)
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Router
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		// TODO: wire up real auth handler for /auth/register and /auth/login

		// Protected routes — using MockAuth until real JWT auth is implemented
		api := v1.Group("/")
		api.Use(middleware.MockAuth())
		{
			api.GET("/profile", userHandler.GetUserByID) // Example protected route for user profile
			api.GET("/persons", personHandler.GetPersons)
			api.POST("/persons", personHandler.CreatePerson)
			api.GET("/persons/:id", personHandler.GetPersonByID)
			api.DELETE("/persons/:id", personHandler.DeletePerson)
		}
	}

	log.Printf("server starting on :%s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
