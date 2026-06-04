package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YuraSahanovskyi/booking-system/internal/db/sqlc"
	"github.com/YuraSahanovskyi/booking-system/internal/handler"
	"github.com/YuraSahanovskyi/booking-system/internal/repository/postgres"
	"github.com/YuraSahanovskyi/booking-system/internal/service"
	"github.com/YuraSahanovskyi/booking-system/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	dbUrl := requireEnv("DB_URL")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	redisUrl := requireEnv("REDIS_URL")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v\n", err)
	}

	queries := sqlc.New(pool)
	userRepo := postgres.NewUserRepository(queries)
	eventRepo := postgres.NewEventRepository(queries)
	bookingRepo := postgres.NewBookingRepository(queries)

	jwtKey := requireEnv("JWT_KEY")
	enableCache := requireBoolEnv("ENABLE_CACHE")

	authService := service.NewAuthService(userRepo, jwtKey, time.Hour*24)
	eventService := service.NewEventService(eventRepo)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, time.Minute*15, rdb, enableCache)

	cleanupWorker := worker.NewCleanupWorker(bookingService, time.Minute)
	go cleanupWorker.Start(ctx)

	handlers := handler.NewHandler(authService, eventService, bookingService)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("invalid PORT")
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handlers.Init(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server started on :" + port)

	<-ctx.Done()

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func requireBoolEnv(key string) bool {
	val := requireEnv(key)
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	log.Fatalf("Environment variable %s must be 'true' or 'false'\n", key)
	return false
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		log.Fatalf("%s is not set\n", key)
	}
	return val
}
