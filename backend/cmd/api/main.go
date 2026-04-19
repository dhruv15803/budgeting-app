package main

import (
	"context"
	"log"
	"net/http"

	"github.com/dhruv15803/budgeting-app/internal/auth"
	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/dhruv15803/budgeting-app/internal/database"
	"github.com/dhruv15803/budgeting-app/internal/email"
	"github.com/dhruv15803/budgeting-app/internal/handlers"
	appmiddleware "github.com/dhruv15803/budgeting-app/internal/middleware"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
	"github.com/dhruv15803/budgeting-app/internal/services"
	"github.com/dhruv15803/budgeting-app/internal/worker"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DbConfig).Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	queue := worker.NewQueue(256)
	sender := email.NewSender(cfg.SMTP)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	worker.RunVerificationMailWorker(ctx, queue, sender.SendVerification, log.Printf)

	jwtSigner := auth.NewJWTSigner(cfg.JWTSecret, cfg.JWTExpiry)

	repo := repositories.NewRepository(db)
	service := services.NewService(repo, cfg, jwtSigner, queue)
	handler := handlers.NewHandler(service)

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {

		r.Get("/health", handler.HealthCheckHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handler.Register)
			r.Post("/login", handler.Login)
			r.Get("/verify-email", handler.VerifyEmail)

			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.AuthMiddleware(jwtSigner))
				r.Get("/me", handler.Me)
			})
		})

		r.Route("/expenses", func(r chi.Router) {
			r.Use(appmiddleware.AuthMiddleware(jwtSigner))
			r.Post("/", handler.CreateExpense)
			r.Get("/", handler.ListExpenses)
			r.Put("/{id}", handler.UpdateExpense)
			r.Delete("/{id}", handler.DeleteExpense)
		})

	})

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	log.Printf("Staring server on %s\n", cfg.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
