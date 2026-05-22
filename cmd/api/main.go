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
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("invalid DB_URL")
	}

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

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatal("invalid REDIS_URL")
	}

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

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("invalid JWT_KEY")
	}

	authService := service.NewAuthService(userRepo, jwtKey, time.Hour*24)
	eventService := service.NewEventService(eventRepo)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, time.Minute*15, rdb)

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
