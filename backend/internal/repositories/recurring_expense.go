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

type RecurringExpenseRepo struct {
	db *sqlx.DB
}

func NewRecurringExpenseRepo(db *sqlx.DB) *RecurringExpenseRepo {
	return &RecurringExpenseRepo{db: db}
}

var recurringSortOrders = map[string]string{
	"next_asc":     "next_occurrence ASC,  id ASC",
	"next_desc":    "next_occurrence DESC, id DESC",
	"created_asc":  "created_at ASC,  id ASC",
	"created_desc": "created_at DESC, id DESC",
}

const recurringColumns = `id, title, description, amount, user_id, category_id, start_date, end_date,
		frequency::text AS frequency, next_occurrence, is_active, created_at, updated_at`

func (r *RecurringExpenseRepo) Create(userID, categoryID int, title string, description *string, amount float64, startDate time.Time, endDate *time.Time, frequency string) (*models.RecurringExpense, error) {
	var out models.RecurringExpense
	err := r.db.QueryRowx(fmt.Sprintf(`
		INSERT INTO recurring_expenses (title, description, amount, user_id, category_id, start_date, end_date, frequency, next_occurrence)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $6)
		RETURNING %s
	`, recurringColumns), title, description, amount, userID, categoryID, startDate, endDate, frequency).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *RecurringExpenseRepo) GetByID(id int) (*models.RecurringExpense, error) {
	var out models.RecurringExpense
	err := r.db.Get(&out, fmt.Sprintf(`SELECT %s FROM recurring_expenses WHERE id = $1`, recurringColumns), id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (r *RecurringExpenseRepo) Update(id int, title string, description *string, amount float64, categoryID int, startDate time.Time, endDate *time.Time, frequency string, isActive bool, nextOccurrence time.Time) (*models.RecurringExpense, error) {
	var out models.RecurringExpense
	err := r.db.QueryRowx(fmt.Sprintf(`
		UPDATE recurring_expenses
		SET title = $2,
		    description = $3,
		    amount = $4,
		    category_id = $5,
		    start_date = $6,
		    end_date = $7,
		    frequency = $8,
		    is_active = $9,
		    next_occurrence = $10,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING %s
	`, recurringColumns), id, title, description, amount, categoryID, startDate, endDate, frequency, isActive, nextOccurrence).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *RecurringExpenseRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM recurring_expenses WHERE id = $1`, id)
	return err
}

func (r *RecurringExpenseRepo) ListByUser(userID int, f models.RecurringExpenseFilter) ([]models.RecurringExpense, int, error) {
	where := []string{"user_id = $1"}
	args := []interface{}{userID}
	n := 2

	if f.Search != nil && *f.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", n, n))
		args = append(args, "%"+*f.Search+"%")
		n++
	}
	if len(f.CategoryIDs) > 0 {
		where = append(where, fmt.Sprintf("category_id = ANY($%d)", n))
		args = append(args, f.CategoryIDs)
		n++
	}
	if f.Frequency != nil && *f.Frequency != "" {
		where = append(where, fmt.Sprintf("frequency = $%d", n))
		args = append(args, *f.Frequency)
		n++
	}
	if f.IsActive != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", n))
		args = append(args, *f.IsActive)
		n++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	orderClause := recurringSortOrders["next_asc"]
	if order, ok := recurringSortOrders[f.SortBy]; ok {
		orderClause = order
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM recurring_expenses %s", whereClause)
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	limit := f.PageSize
	offset := (f.Page - 1) * f.PageSize

	selectArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT %s
		FROM recurring_expenses
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, recurringColumns, whereClause, orderClause, n, n+1)

	var rows []models.RecurringExpense
	if err := r.db.Select(&rows, dataQuery, selectArgs...); err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *RecurringExpenseRepo) ListDue(today time.Time) ([]models.RecurringExpense, error) {
	var rows []models.RecurringExpense
	err := r.db.Select(&rows, fmt.Sprintf(`
		SELECT %s
		FROM recurring_expenses
		WHERE is_active = TRUE
		  AND next_occurrence <= $1
		ORDER BY id ASC
	`, recurringColumns), today)
	return rows, err
}

func (r *RecurringExpenseRepo) AdvanceTx(tx *sqlx.Tx, id int, newNext time.Time, deactivate bool) error {
	_, err := tx.Exec(`
		UPDATE recurring_expenses
		SET next_occurrence = $2,
		    is_active = CASE WHEN $3 THEN FALSE ELSE is_active END,
		    updated_at = NOW()
		WHERE id = $1
	`, id, newNext, deactivate)
	return err
}
