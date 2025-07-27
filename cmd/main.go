package main

import (
	"context"
	"fmt"
	"log"
	"my-pastebin/internal/api"
	"my-pastebin/internal/paste"
	"my-pastebin/internal/storage"
	"my-pastebin/internal/tracing"
	"os"
	"time"

	_ "my-pastebin/docs"
	"my-pastebin/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Важно: импортируем сгенерированную документацию

// @title           my-pastebin Service API
// @version         1.0
// @description     A minimalist service for sharing text snippets.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	tracerProvider, shutdown := tracing.InitTracerProvider("jaeger:4317", "pastebin-service")
	defer shutdown()

	otel.SetTracerProvider(tracerProvider)

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

	appMetrics := metrics.NewMetrics(prometheus.DefaultRegisterer)
	dbStorage := storage.New(db, appMetrics)
	apiHandler := api.New(dbStorage, appMetrics)

	go startCleanupWorker(dbStorage)

	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(appMetrics.PrometheusMiddleware())
	router.Use(otelgin.Middleware("pastebin-service", otelgin.WithTracerProvider(tracerProvider)))

	apiHandler.RegisterRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", metrics.PrometheusHandler())

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
