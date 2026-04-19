package repositories

import (
	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type ExpenseCategoryRepo struct {
	db *sqlx.DB
}

func NewExpenseCategoryRepo(db *sqlx.DB) *ExpenseCategoryRepo {
	return &ExpenseCategoryRepo{db: db}
}

func (r *ExpenseCategoryRepo) ListAll() ([]models.ExpenseCategory, error) {
	var rows []models.ExpenseCategory
	err := r.db.Select(&rows, `
		SELECT id, category_name, created_at, updated_at
		FROM expense_categories
		ORDER BY category_name ASC
	`)
	return rows, err
}
