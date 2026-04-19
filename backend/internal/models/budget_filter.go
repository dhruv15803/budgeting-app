package models

type BudgetFilter struct {
	// optional: restrict to budgets in a given calendar year
	Year *int

	// "month_asc" (default) | "month_desc"
	SortBy string

	// pagination
	Page     int
	PageSize int
}
