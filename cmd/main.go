package main

import (
	"context"
	"fmt"
	"log"
	"my-pastebin/internal/api"
	"my-pastebin/internal/paste"
	"my-pastebin/internal/storage"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Running DB migrations...")
	if err := db.AutoMigrate(&paste.Paste{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	dbStorage := storage.New(db)
	apiHandler := api.New(dbStorage)

	go startCleanupWorker(dbStorage)

	router := gin.Default()
	apiHandler.RegisterRoutes(router)

	serverPort := ":" + os.Getenv("SERVER_PORT")
	log.Printf("Starting server on port %s", serverPort)
	if err := router.Run(serverPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func startCleanupWorker(s *storage.Storage) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		log.Println("Running cleanup job for expired pastes...")
		deletedCount, err := s.DeleteExpired(context.Background())
		if err != nil {
			log.Printf("Error cleaning up expired pastes: %v", err)
		} else {
			log.Printf("Cleaned up %d expired pastes", deletedCount)
		}
	}
}
