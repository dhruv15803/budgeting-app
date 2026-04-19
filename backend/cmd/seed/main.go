package main

import (
	"log"

	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/dhruv15803/budgeting-app/internal/database"
)

var expenseCategories = []string{
	"Housing",
	"Food & Dining",
	"Transportation",
	"Utilities",
	"Healthcare",
	"Entertainment",
	"Shopping",
	"Education",
	"Travel",
	"Personal Care",
	"Insurance",
	"Investments & Savings",
	"Subscriptions",
	"Gifts & Donations",
	"Childcare",
	"Pet Care",
	"Fitness & Sports",
	"Technology",
	"Taxes",
	"Miscellaneous",
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DbConfig).Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database. Seeding expense categories...")

	inserted := 0
	skipped := 0

	for _, name := range expenseCategories {
		res, err := db.Exec(
			`INSERT INTO expense_categories (category_name) VALUES ($1) ON CONFLICT (category_name) DO NOTHING`,
			name,
		)
		if err != nil {
			log.Fatalf("failed to insert category %q: %v", name, err)
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			skipped++
			log.Printf("  skipped (already exists): %s", name)
		} else {
			inserted++
			log.Printf("  inserted: %s", name)
		}
	}

	log.Printf("Done. %d inserted, %d skipped.", inserted, skipped)
}
