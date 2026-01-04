package main

import (
	"context"
	"embed"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/Greezaaa/auth-go/internal/user"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	_ = godotenv.Load()
	logger := log.Default()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.New(context.Background(), dbURL)
		if err == nil && pool.Ping(context.Background()) == nil {
			break
		}
		logger.Printf("Waiting for DB (attempt %d/10)...", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logger.Fatalf("DB Connection failed: %v", err)
	}
	defer pool.Close()

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatalf("Goose dialect error: %v", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		logger.Fatalf("Migration failed: %v", err)
	}
	db.Close()
	logger.Println("✅ Migrations applied")

	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo, logger)
	userHandler := user.NewHandler(userService, logger)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiPath := fmt.Sprintf("/%s/%s", os.Getenv("API_PREFIX"), os.Getenv("API_VERSION"))
	r.Route(apiPath, func(r chi.Router) {
		r.Mount("/users", userHandler.Routes())
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Printf("🚀 Live at: http://localhost:%s%s", port, apiPath)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
