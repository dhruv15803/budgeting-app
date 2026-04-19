package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type ExpenseRepo struct {
	db *sqlx.DB
}

func NewExpenseRepo(db *sqlx.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (e *ExpenseRepo) Create(userID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error) {
	var out models.Expense
	err := e.db.QueryRowx(`
		INSERT INTO expenses (title, description, amount, user_id, category_id, expense_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, amount, user_id, category_id, recurring_expense_id, expense_date, created_at, updated_at
	`, title, description, amount, userID, categoryID, expenseDate).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (e *ExpenseRepo) GetByID(id int) (*models.Expense, error) {
	var out models.Expense
	err := e.db.Get(&out, `
		SELECT id, title, description, amount, user_id, category_id, recurring_expense_id, expense_date, created_at, updated_at
		FROM expenses WHERE id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (e *ExpenseRepo) Update(id int, title string, description *string, amount float64, categoryID int, expenseDate time.Time) (*models.Expense, error) {
	var out models.Expense
	err := e.db.QueryRowx(`
		UPDATE expenses
		SET title = $2, description = $3, amount = $4, category_id = $5, expense_date = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, amount, user_id, category_id, recurring_expense_id, expense_date, created_at, updated_at
	`, id, title, description, amount, categoryID, expenseDate).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (e *ExpenseRepo) Delete(id int) error {
	_, err := e.db.Exec(`DELETE FROM expenses WHERE id = $1`, id)
	return err
}

func (e *ExpenseRepo) ListByUser(userID, limit, offset int) ([]models.Expense, int, error) {
	var total int
	if err := e.db.Get(&total, `SELECT COUNT(*) FROM expenses WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	var rows []models.Expense
	if err := e.db.Select(&rows, `
		SELECT id, title, description, amount, user_id, category_id, recurring_expense_id, expense_date, created_at, updated_at
		FROM expenses
		WHERE user_id = $1
		ORDER BY expense_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset); err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
