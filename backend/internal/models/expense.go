package models

import "time"

type ExpenseCategory struct {
	ID           int        `db:"id" json:"id"`
	CategoryName string     `db:"category_name" json:"category_name"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at"`
}

type Expense struct {
	ID                 int        `db:"id"`
	Title              string     `db:"title"`
	Description        *string    `db:"description"`
	Amount             float64    `db:"amount"`
	UserID             int        `db:"user_id"`
	CategoryID         int        `db:"category_id"`
	RecurringExpenseID *int       `db:"recurring_expense_id"`
	ExpenseDate        time.Time  `db:"expense_date"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          *time.Time `db:"updated_at"`
}
