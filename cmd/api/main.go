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
	dbUrl := "postgres://postgres:postgres@db:5432/postgres"

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

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v\n", err)
	}

	queries := sqlc.New(pool)
	userRepo := postgres.NewUserRepository(queries)
	eventRepo := postgres.NewEventRepository(queries)
	bookingRepo := postgres.NewBookingRepository(queries)

	authService := service.NewAuthService(userRepo, "your-super-secret-key", time.Hour*24)
	eventService := service.NewEventService(eventRepo)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, time.Minute*15, rdb)

	cleanupWorker := worker.NewCleanupWorker(bookingService, time.Minute)
	go cleanupWorker.Start(ctx)

	handlers := handler.NewHandler(authService, eventService, bookingService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handlers.Init(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server started on :8080")

	<-ctx.Done()

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
