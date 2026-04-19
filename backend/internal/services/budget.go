package services

import (
	"fmt"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
)

type BudgetServiceImpl struct {
	repo *repositories.Repository
}

func NewBudgetService(repo *repositories.Repository) *BudgetServiceImpl {
	return &BudgetServiceImpl{repo: repo}
}

// BudgetOverviewResult is the assembled response for GET /budgets/{month}.
type BudgetOverviewResult struct {
	Budget     *models.MonthlyBudget
	TotalSpent float64
	Remaining  float64
	Categories []models.CategoryBudgetWithSpending
}

// monthRange returns the first and last day of the month that budgetMonth falls in.
func monthRange(budgetMonth time.Time) (from, to time.Time) {
	first := time.Date(budgetMonth.Year(), budgetMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1)
	return first, last
}

func (s *BudgetServiceImpl) getBudgetByMonth(userID int, budgetMonth time.Time) (*models.MonthlyBudget, error) {
	budget, err := s.repo.Budgets.GetBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return nil, err
	}
	if budget == nil {
		return nil, ErrNotFound
	}
	return budget, nil
}

func (s *BudgetServiceImpl) CreateBudget(userID int, budgetMonth time.Time, totalAmount float64) (*models.MonthlyBudget, error) {
	if totalAmount < 0 {
		return nil, fmt.Errorf("total_amount must be 0 or greater")
	}

	existing, err := s.repo.Budgets.GetBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("a budget for %s already exists", budgetMonth.Format("2006-01"))
	}

	return s.repo.Budgets.CreateBudget(userID, budgetMonth, totalAmount)
}

func (s *BudgetServiceImpl) GetBudget(userID int, budgetMonth time.Time) (*BudgetOverviewResult, error) {
	budget, err := s.getBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return nil, err
	}

	from, to := monthRange(budgetMonth)

	totalSpent, err := s.repo.Budgets.GetTotalSpending(userID, from, to)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.Budgets.GetCategorySpending(userID, budget.ID, from, to)
	if err != nil {
		return nil, err
	}

	return &BudgetOverviewResult{
		Budget:     budget,
		TotalSpent: totalSpent,
		Remaining:  budget.TotalAmount - totalSpent,
		Categories: categories,
	}, nil
}

func (s *BudgetServiceImpl) UpdateBudgetTotal(userID int, budgetMonth time.Time, totalAmount float64) (*models.MonthlyBudget, error) {
	if totalAmount < 0 {
		return nil, fmt.Errorf("total_amount must be 0 or greater")
	}

	budget, err := s.getBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return nil, err
	}

	// Ensure the new total is not below current sum of category allocations.
	sum, err := s.repo.Budgets.SumCategoryAllocations(budget.ID)
	if err != nil {
		return nil, err
	}
	if totalAmount < sum {
		return nil, fmt.Errorf("total_amount (%.2f) cannot be less than the sum of category allocations (%.2f)", totalAmount, sum)
	}

	return s.repo.Budgets.UpdateBudgetTotal(budget.ID, totalAmount)
}

func (s *BudgetServiceImpl) DeleteBudget(userID int, budgetMonth time.Time) error {
	budget, err := s.getBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return err
	}
	return s.repo.Budgets.DeleteBudget(budget.ID)
}

func (s *BudgetServiceImpl) ListBudgets(userID int, f models.BudgetFilter) ([]models.MonthlyBudget, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	return s.repo.Budgets.ListBudgets(userID, f)
}

func (s *BudgetServiceImpl) UpsertCategoryBudget(userID int, budgetMonth time.Time, categoryID int, amount float64) (*models.MonthlyBudgetCategory, error) {
	if amount < 0 {
		return nil, fmt.Errorf("allocated_amount must be 0 or greater")
	}
	if categoryID <= 0 {
		return nil, fmt.Errorf("category_id is required")
	}

	budget, err := s.getBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return nil, err
	}

	// Check that adding/updating this allocation does not push total over budget.
	sum, err := s.repo.Budgets.SumCategoryAllocations(budget.ID)
	if err != nil {
		return nil, err
	}

	// If this category already has an allocation, subtract its current value first.
	existing, err := s.repo.Budgets.ListCategoryBudgets(budget.ID)
	if err != nil {
		return nil, err
	}
	var currentAlloc float64
	for _, c := range existing {
		if c.CategoryID == categoryID {
			currentAlloc = c.AllocatedAmount
			break
		}
	}
	newSum := sum - currentAlloc + amount
	if newSum > budget.TotalAmount {
		return nil, fmt.Errorf("total category allocations (%.2f) would exceed total budget (%.2f)", newSum, budget.TotalAmount)
	}

	return s.repo.Budgets.UpsertCategoryBudget(budget.ID, categoryID, amount)
}

func (s *BudgetServiceImpl) DeleteCategoryBudget(userID int, budgetMonth time.Time, categoryID int) error {
	budget, err := s.getBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return err
	}
	if err := s.repo.Budgets.DeleteCategoryBudget(budget.ID, categoryID); err != nil {
		return ErrNotFound
	}
	return nil
}

func (s *BudgetServiceImpl) BulkSetCategoryBudgets(userID int, budgetMonth time.Time, items []models.BulkCategoryItem) error {
	budget, err := s.getBudgetByMonth(userID, budgetMonth)
	if err != nil {
		return err
	}

	var totalAlloc float64
	seen := make(map[int]struct{})
	for _, item := range items {
		if item.AllocatedAmount < 0 {
			return fmt.Errorf("allocated_amount must be 0 or greater for all categories")
		}
		if item.CategoryID <= 0 {
			return fmt.Errorf("all category_id values must be positive")
		}
		if _, dup := seen[item.CategoryID]; dup {
			return fmt.Errorf("duplicate category_id %d in bulk set", item.CategoryID)
		}
		seen[item.CategoryID] = struct{}{}
		totalAlloc += item.AllocatedAmount
	}

	if totalAlloc > budget.TotalAmount {
		return fmt.Errorf("total category allocations (%.2f) would exceed total budget (%.2f)", totalAlloc, budget.TotalAmount)
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.repo.Budgets.BulkSetCategoryBudgetsTx(tx, budget.ID, items); err != nil {
		return err
	}
	return tx.Commit()
}
