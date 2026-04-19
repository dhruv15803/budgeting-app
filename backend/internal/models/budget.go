package models

import "time"

type MonthlyBudget struct {
	ID          int        `db:"id"`
	UserID      int        `db:"user_id"`
	BudgetMonth time.Time  `db:"budget_month"`
	TotalAmount float64    `db:"total_amount"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type MonthlyBudgetCategory struct {
	ID              int        `db:"id"`
	MonthlyBudgetID int        `db:"monthly_budget_id"`
	CategoryID      int        `db:"category_id"`
	AllocatedAmount float64    `db:"allocated_amount"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

// CategoryBudgetWithSpending is a query result type used only by the budget overview.
type CategoryBudgetWithSpending struct {
	ID              int     `db:"id"`
	CategoryID      int     `db:"category_id"`
	CategoryName    string  `db:"category_name"`
	AllocatedAmount float64 `db:"allocated_amount"`
	SpentAmount     float64 `db:"spent_amount"`
}

// BulkCategoryItem is the input type for BulkSetCategoryBudgets.
type BulkCategoryItem struct {
	CategoryID      int
	AllocatedAmount float64
}
