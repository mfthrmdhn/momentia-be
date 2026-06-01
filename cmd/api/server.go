package main

import (
	"log"
	"momentia-be/endpoints"
	"momentia-be/internal/config"
	"momentia-be/internal/handler"
	"momentia-be/repository"
	"momentia-be/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("database connection failed: " + err.Error())
	}
	runMigrations(cfg.DB.DSN())

	// Wire up dependencies
	personRepo := repository.NewPersonRepository(db)
	personHandler := handler.NewPersonHandler(personRepo)
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userRepo, userService)

	// Router
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

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

	//endpoints
	endpoints.RegisterPersonRoutes(r, *personHandler)
	endpoints.RegisterUserRoutes(r, *userHandler)

	log.Printf("server starting on :%s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
