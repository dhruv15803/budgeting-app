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

type BudgetRepo struct {
	db *sqlx.DB
}

func NewBudgetRepo(db *sqlx.DB) *BudgetRepo {
	return &BudgetRepo{db: db}
}

const budgetColumns = `id, user_id, budget_month, total_amount::float8 AS total_amount, created_at, updated_at`
const catBudgetColumns = `id, monthly_budget_id, category_id, allocated_amount::float8 AS allocated_amount, created_at, updated_at`

var budgetSortOrders = map[string]string{
	"month_asc":  "budget_month ASC",
	"month_desc": "budget_month DESC",
}

// ── monthly_budgets ──────────────────────────────────────────────────────────

func (r *BudgetRepo) CreateBudget(userID int, budgetMonth time.Time, totalAmount float64) (*models.MonthlyBudget, error) {
	var out models.MonthlyBudget
	err := r.db.QueryRowx(fmt.Sprintf(`
		INSERT INTO monthly_budgets (user_id, budget_month, total_amount)
		VALUES ($1, $2, $3)
		RETURNING %s
	`, budgetColumns), userID, budgetMonth, totalAmount).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *BudgetRepo) GetBudgetByMonth(userID int, budgetMonth time.Time) (*models.MonthlyBudget, error) {
	var out models.MonthlyBudget
	err := r.db.Get(&out, fmt.Sprintf(`
		SELECT %s FROM monthly_budgets WHERE user_id = $1 AND budget_month = $2
	`, budgetColumns), userID, budgetMonth)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (r *BudgetRepo) GetBudgetByID(id int) (*models.MonthlyBudget, error) {
	var out models.MonthlyBudget
	err := r.db.Get(&out, fmt.Sprintf(`
		SELECT %s FROM monthly_budgets WHERE id = $1
	`, budgetColumns), id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (r *BudgetRepo) UpdateBudgetTotal(id int, totalAmount float64) (*models.MonthlyBudget, error) {
	var out models.MonthlyBudget
	err := r.db.QueryRowx(fmt.Sprintf(`
		UPDATE monthly_budgets
		SET total_amount = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING %s
	`, budgetColumns), id, totalAmount).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *BudgetRepo) DeleteBudget(id int) error {
	_, err := r.db.Exec(`DELETE FROM monthly_budgets WHERE id = $1`, id)
	return err
}

func (r *BudgetRepo) ListBudgets(userID int, f models.BudgetFilter) ([]models.MonthlyBudget, int, error) {
	where := []string{"user_id = $1"}
	args := []interface{}{userID}
	n := 2

	if f.Year != nil {
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM budget_month) = $%d", n))
		args = append(args, *f.Year)
		n++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	orderClause := budgetSortOrders["month_asc"]
	if order, ok := budgetSortOrders[f.SortBy]; ok {
		orderClause = order
	}

	var total int
	if err := r.db.Get(&total, fmt.Sprintf("SELECT COUNT(*) FROM monthly_budgets %s", whereClause), args...); err != nil {
		return nil, 0, err
	}

	limit := f.PageSize
	offset := (f.Page - 1) * f.PageSize
	selectArgs := append(args, limit, offset)

	dataQuery := fmt.Sprintf(`
		SELECT %s FROM monthly_budgets
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, budgetColumns, whereClause, orderClause, n, n+1)

	var rows []models.MonthlyBudget
	if err := r.db.Select(&rows, dataQuery, selectArgs...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ── monthly_category_budgets ─────────────────────────────────────────────────

func (r *BudgetRepo) UpsertCategoryBudget(budgetID, categoryID int, amount float64) (*models.MonthlyBudgetCategory, error) {
	var out models.MonthlyBudgetCategory
	err := r.db.QueryRowx(fmt.Sprintf(`
		INSERT INTO monthly_category_budgets (monthly_budget_id, category_id, allocated_amount)
		VALUES ($1, $2, $3)
		ON CONFLICT (monthly_budget_id, category_id)
		DO UPDATE SET allocated_amount = EXCLUDED.allocated_amount, updated_at = NOW()
		RETURNING %s
	`, catBudgetColumns), budgetID, categoryID, amount).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *BudgetRepo) DeleteCategoryBudget(budgetID, categoryID int) error {
	res, err := r.db.Exec(`DELETE FROM monthly_category_budgets WHERE monthly_budget_id = $1 AND category_id = $2`, budgetID, categoryID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *BudgetRepo) ListCategoryBudgets(budgetID int) ([]models.MonthlyBudgetCategory, error) {
	var rows []models.MonthlyBudgetCategory
	err := r.db.Select(&rows, fmt.Sprintf(`
		SELECT %s FROM monthly_category_budgets
		WHERE monthly_budget_id = $1
		ORDER BY category_id ASC
	`, catBudgetColumns), budgetID)
	return rows, err
}

func (r *BudgetRepo) SumCategoryAllocations(budgetID int) (float64, error) {
	var sum float64
	err := r.db.Get(&sum, `
		SELECT COALESCE(SUM(allocated_amount), 0)::float8 FROM monthly_category_budgets WHERE monthly_budget_id = $1
	`, budgetID)
	return sum, err
}

// BulkSetCategoryBudgetsTx deletes all existing category allocations for budgetID and inserts
// the new set, all within the provided transaction.
func (r *BudgetRepo) BulkSetCategoryBudgetsTx(tx *sqlx.Tx, budgetID int, items []models.BulkCategoryItem) error {
	if _, err := tx.Exec(`DELETE FROM monthly_category_budgets WHERE monthly_budget_id = $1`, budgetID); err != nil {
		return err
	}
	for _, item := range items {
		if _, err := tx.Exec(`
			INSERT INTO monthly_category_budgets (monthly_budget_id, category_id, allocated_amount)
			VALUES ($1, $2, $3)
		`, budgetID, item.CategoryID, item.AllocatedAmount); err != nil {
			return err
		}
	}
	return nil
}

// ── overview queries ─────────────────────────────────────────────────────────

// GetCategorySpending returns each category allocation for the budget joined with
// actual expense totals for the given date range.
func (r *BudgetRepo) GetCategorySpending(userID, budgetID int, from, to time.Time) ([]models.CategoryBudgetWithSpending, error) {
	var rows []models.CategoryBudgetWithSpending
	err := r.db.Select(&rows, `
		SELECT
		    mcb.id,
		    mcb.category_id,
		    ec.category_name,
		    mcb.allocated_amount::float8 AS allocated_amount,
		    COALESCE(SUM(e.amount), 0)::float8   AS spent_amount
		FROM monthly_category_budgets mcb
		JOIN expense_categories ec ON ec.id = mcb.category_id
		LEFT JOIN expenses e
		    ON  e.category_id  = mcb.category_id
		    AND e.user_id      = $1
		    AND e.expense_date >= $3
		    AND e.expense_date <= $4
		WHERE mcb.monthly_budget_id = $2
		GROUP BY mcb.id, mcb.category_id, ec.category_name, mcb.allocated_amount
		ORDER BY ec.category_name ASC
	`, userID, budgetID, from, to)
	return rows, err
}

// GetTotalSpending returns the sum of expenses for the user in the given date range.
func (r *BudgetRepo) GetTotalSpending(userID int, from, to time.Time) (float64, error) {
	var total float64
	err := r.db.Get(&total, `
		SELECT COALESCE(SUM(amount), 0)::float8
		FROM expenses
		WHERE user_id = $1
		  AND expense_date >= $2
		  AND expense_date <= $3
	`, userID, from, to)
	return total, err
}
