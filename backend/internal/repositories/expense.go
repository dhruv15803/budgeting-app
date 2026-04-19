package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

// CreateFromRecurringTx inserts an expense derived from a recurring expense on the given date.
// Uses ON CONFLICT on the partial unique index (recurring_expense_id, expense_date) so repeated
// generator runs are safe and never produce duplicates. Returns true if a row was actually inserted.
func (e *ExpenseRepo) CreateFromRecurringTx(tx *sqlx.Tx, r *models.RecurringExpense, date time.Time) (bool, error) {
	res, err := tx.Exec(`
		INSERT INTO expenses (title, description, amount, user_id, category_id, recurring_expense_id, expense_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (recurring_expense_id, expense_date)
		WHERE recurring_expense_id IS NOT NULL
		DO NOTHING
	`, r.Title, r.Description, r.Amount, r.UserID, r.CategoryID, r.ID, date)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

var sortOrders = map[string]string{
	"date_asc":    "expense_date ASC,  created_at ASC",
	"date_desc":   "expense_date DESC, created_at DESC",
	"amount_asc":  "amount ASC,  expense_date DESC",
	"amount_desc": "amount DESC, expense_date DESC",
}

func (e *ExpenseRepo) ListByUser(userID int, f models.ExpenseFilter) ([]models.Expense, int, error) {
	where := []string{"user_id = $1"}
	args := []interface{}{userID}
	n := 2

	if f.Search != nil && *f.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", n, n))
		args = append(args, "%"+*f.Search+"%")
		n++
	}
	if f.DateFrom != nil {
		where = append(where, fmt.Sprintf("expense_date >= $%d", n))
		args = append(args, *f.DateFrom)
		n++
	}
	if f.DateTo != nil {
		where = append(where, fmt.Sprintf("expense_date <= $%d", n))
		args = append(args, *f.DateTo)
		n++
	}
	if len(f.CategoryIDs) > 0 {
		where = append(where, fmt.Sprintf("category_id = ANY($%d)", n))
		args = append(args, f.CategoryIDs)
		n++
	}
	if f.AmountMin != nil {
		where = append(where, fmt.Sprintf("amount >= $%d", n))
		args = append(args, *f.AmountMin)
		n++
	}
	if f.AmountMax != nil {
		where = append(where, fmt.Sprintf("amount <= $%d", n))
		args = append(args, *f.AmountMax)
		n++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	orderClause := sortOrders["date_desc"]
	if order, ok := sortOrders[f.SortBy]; ok {
		orderClause = order
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses %s", whereClause)
	if err := e.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	limit := f.PageSize
	offset := (f.Page - 1) * f.PageSize

	selectArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT id, title, description, amount, user_id, category_id, recurring_expense_id, expense_date, created_at, updated_at
		FROM expenses
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, n, n+1)

	var rows []models.Expense
	if err := e.db.Select(&rows, dataQuery, selectArgs...); err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
