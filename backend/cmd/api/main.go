package main

import (
	"context"
	"log"
	"net/http"
	"time"

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
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	service := services.NewService(repo, cfg, jwtSigner, queue, log.Printf)
	handler := handlers.NewHandler(service)

	scheduler, err := worker.NewRecurringScheduler(cfg.CronSchedule, service.RecurringExpenses.RunDueGenerator, log.Printf)
	if err != nil {
		log.Fatalf("Error initialising recurring scheduler: %v", err)
	}
	scheduler.Start()
	defer scheduler.Stop()

	// Run once on startup to catch up any missed occurrences during downtime.
	go func() {
		startupCtx, cancelStartup := context.WithTimeout(ctx, 5*time.Minute)
		defer cancelStartup()
		if n, err := scheduler.RunNow(startupCtx); err != nil {
			log.Printf("startup recurring catch-up failed: %v", err)
		} else {
			log.Printf("startup recurring catch-up created %d expenses", n)
		}
	}()

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {

		r.Get("/health", handler.HealthCheckHandler)
		r.Get("/expense-categories", handler.ListExpenseCategories)

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

		r.Route("/recurring-expenses", func(r chi.Router) {
			r.Use(appmiddleware.AuthMiddleware(jwtSigner))
			r.Post("/", handler.CreateRecurringExpense)
			r.Get("/", handler.ListRecurringExpenses)
			r.Get("/{id}", handler.GetRecurringExpense)
			r.Put("/{id}", handler.UpdateRecurringExpense)
			r.Delete("/{id}", handler.DeleteRecurringExpense)
		})

		r.Route("/budgets", func(r chi.Router) {
			r.Use(appmiddleware.AuthMiddleware(jwtSigner))
			r.Post("/", handler.CreateBudget)
			r.Get("/", handler.ListBudgets)
			r.Get("/{month}", handler.GetBudget)
			r.Put("/{month}", handler.UpdateBudget)
			r.Delete("/{month}", handler.DeleteBudget)
			r.Post("/{month}/categories", handler.UpsertCategoryBudget)
			r.Delete("/{month}/categories/{category_id}", handler.DeleteCategoryBudget)
			r.Put("/{month}/categories", handler.BulkSetCategoryBudgets)
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
