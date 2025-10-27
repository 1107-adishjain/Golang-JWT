package main

import (
	"context"
	"github.com/1107-adishjain/golang-jwt/internal/config"
	"github.com/1107-adishjain/golang-jwt/internal/database"
	"github.com/1107-adishjain/golang-jwt/internal/models"
	"github.com/1107-adishjain/golang-jwt/internal/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App struct holds the router, database connection, and configuration
type App struct {
	Router *gin.Engine
	DB     *gorm.DB
	Config *config.Config
}

func main() {
	cfg := config.LoadConfig()

	db, err := database.DBinitialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	} else {
		log.Printf("successfully connected to database")
	}
	db.AutoMigrate(&models.User{})
	// Ensure DB is closed on exit.
	defer func() {
		if err := database.DBClose(db); err != nil {
			log.Printf("error closing db: %v", err)
		} else {
			log.Printf("database connection closed")
		}
	}()

	router := gin.Default()

	app := &App{
		Router: router,
		DB:     db,
		Config: cfg,
	}

	// Basic health route
	app.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Register application routes
	routes.AuthRoutes(app.Router, app.DB)
	routes.UserRoutes(app.Router, app.DB)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine so we can listen for shutdown signals.
	go func() {
		log.Printf("starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Listen for interrupt signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	} else {
		log.Println("server stopped gracefully")
	}
}
