package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/viditagrawal56/url-shortner/internal/config"
	"github.com/viditagrawal56/url-shortner/internal/db"
	"github.com/viditagrawal56/url-shortner/internal/handlers"
	"github.com/viditagrawal56/url-shortner/internal/models"
)

func main() {
	// Load the config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	resetDB := flag.Bool("reset-db", false, "Reset the database before running migrations")
	flag.Parse()

	// Load and open a connection to the database
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Connected to the database")
	defer database.Close()

	// Models to migrate
	models := []any{
		&models.User{},
		&models.AuthorizedEmail{},
		&models.ShortURL{},
		&models.Credentials{},
		&models.TemporaryToken{},
	}

	// Handle migrations based on flag
	if *resetDB {
		if err := database.ResetAndMigrate(models...); err != nil {
			log.Fatalf("Failed to reset and migrate database: %v", err)
		}
		log.Println("Database reset and migration completed")
	} else {
		if err := database.AutoMigrate(models...); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database migration completed")
	}

	// Initialize router with handlers
	router := handlers.NewRouter(database, cfg)

	// Configure the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//Start the server
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Wait for the current operations to compelete
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
