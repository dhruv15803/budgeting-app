package models

type RecurringExpenseFilter struct {
	// ILIKE on title OR description
	Search *string

	// repeatable ?category_id=1&category_id=3
	CategoryIDs []int

	// daily | weekly | monthly | yearly
	Frequency *string

	// true/false
	IsActive *bool

	// sort_by: "next_asc" (default) | "next_desc" | "created_asc" | "created_desc"
	SortBy string

	// pagination
	Page     int
	PageSize int
}
