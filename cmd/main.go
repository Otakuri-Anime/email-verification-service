package main

import (
	"context"
	"email-verification/internal/config"
	"email-verification/internal/delivery/handler"
	"email-verification/internal/email"
	"email-verification/internal/repository"
	"email-verification/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Проверка подключения к Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	verificationRepo := repository.NewRedisVerificationRepository(redisClient)

	// Инициализация email сервиса
	emailSender, err := email.NewElasticEmailSender(
		cfg.Email.APIKey,
		cfg.Email.FromEmail,
		cfg.Email.Endpoint,
	)
	if err != nil {
		log.Fatalf("Failed to create email sender: %v", err)
	}

	// Инициализация сервиса
	verificationSvc := service.NewVerificationService(
		verificationRepo,
		emailSender,
		cfg.Email.CodeLength,
		cfg.Email.CodeExpiry,
	)

	// Инициализация HTTP обработчиков
	handler := handler.NewVerificationHandler(verificationSvc, 10*time.Second)

	// Настройка маршрутов
	mux := http.NewServeMux()
	handler.SetupRoutes(mux)

	// Настройка CORS
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Создание HTTP сервера
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: corsMiddleware(mux),
	}

	// Graceful shutdown
	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server started on port %s", cfg.Server.Port)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received %v signal, shutting down...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Fatalf("Force shutdown failed: %v", err)
			}
		}
	}

	log.Println("Server stopped")
}
