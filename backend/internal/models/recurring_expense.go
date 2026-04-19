package models

import "time"

const (
	FrequencyDaily   = "daily"
	FrequencyWeekly  = "weekly"
	FrequencyMonthly = "monthly"
	FrequencyYearly  = "yearly"
)

type RecurringExpense struct {
	ID             int        `db:"id"`
	Title          string     `db:"title"`
	Description    *string    `db:"description"`
	Amount         float64    `db:"amount"`
	UserID         int        `db:"user_id"`
	CategoryID     int        `db:"category_id"`
	StartDate      time.Time  `db:"start_date"`
	EndDate        *time.Time `db:"end_date"`
	Frequency      string     `db:"frequency"`
	NextOccurrence time.Time  `db:"next_occurrence"`
	IsActive       bool       `db:"is_active"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}
