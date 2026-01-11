package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"dkmbackend/internal/config"
	"dkmbackend/internal/db"
	"dkmbackend/internal/httpserver"
	"dkmbackend/internal/repository/mongoimpl"
	"dkmbackend/internal/services"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := db.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	defer func() {
		_ = mongoClient.Disconnect(context.Background())
	}()

	database := mongoClient.Database(cfg.MongoDB)

	// repositories
	blogRepo := mongoimpl.NewBlogRepository(database)
	productRepo := mongoimpl.NewProductRepository(database)
	careerRepo := mongoimpl.NewCareerRepository(database)

	// services
	blogSvc := services.NewBlogService(blogRepo)
	productSvc := services.NewProductService(productRepo)
	careerSvc := services.NewCareerService(careerRepo)
	authSvc := services.NewAuthService(cfg)

	// http server
	router := httpserver.NewRouter(cfg, blogSvc, productSvc, careerSvc, authSvc)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// graceful shutdown
	go func() {
		log.Printf("server listening on :%d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
