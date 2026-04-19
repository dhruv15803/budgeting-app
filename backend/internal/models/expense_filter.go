package models

import "time"

type ExpenseFilter struct {
	// full-text search across title and description (ILIKE)
	Search *string

	// date range — both optional, combinable freely
	DateFrom *time.Time // expense_date >= DateFrom
	DateTo   *time.Time // expense_date <= DateTo

	// category filter — multiple IDs allowed
	CategoryIDs []int

	// amount range
	AmountMin *float64
	AmountMax *float64

	// sort_by: "date_asc" | "date_desc" (default) | "amount_asc" | "amount_desc"
	SortBy string

	// pagination
	Page     int
	PageSize int
}
