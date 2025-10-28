package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bowe99/phone-usage-service/internal/api/handler"
	"github.com/bowe99/phone-usage-service/internal/api/router"
	"github.com/bowe99/phone-usage-service/internal/application/service"
	"github.com/bowe99/phone-usage-service/internal/infra/config"
	"github.com/bowe99/phone-usage-service/internal/infra/database"
	"github.com/bowe99/phone-usage-service/internal/infra/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Connecting to MongoDB...")
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database, cfg.MongoDB.Timeout)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()
	log.Println("Successfully connected to MongoDB")

	userRepo := repository.SetupUserRepository(db.Database)
	// cycleRepo := repository.NewCycleRepository(db.Database)
	// usageRepo := repository.NewUsageRepository(db.Database)

	// Initialize services (Application layer)
	userService := service.SetupUserService(userRepo)
	// cycleService := service.NewCycleService(cycleRepo)
	// usageService := service.NewUsageService(usageRepo, cycleRepo)

	// Initialize handlers (Presentation layer)
	userHandler := handler.SetupUserHandler(userService)
	// cycleHandler := handler.NewCycleHandler(cycleService)
	// usageHandler := handler.NewUsageHandler(usageService)

	r := setupRouter(db, cfg, userHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s...", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(db *database.MongoDB, cfg *config.Config, userHandler *handler.UserHandler) *gin.Engine {
	return router.SetupRouter(db, cfg.Server.GinMode, userHandler)
}