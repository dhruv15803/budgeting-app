package services

import (
	"fmt"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
)

type ExpenseServiceImpl struct {
	repo *repositories.Repository
}

func NewExpenseService(repo *repositories.Repository) *ExpenseServiceImpl {
	return &ExpenseServiceImpl{repo: repo}
}

func (s *ExpenseServiceImpl) CreateExpense(userID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if categoryID <= 0 {
		return nil, fmt.Errorf("category_id is required")
	}
	if expenseDate.IsZero() {
		return nil, fmt.Errorf("expense_date is required")
	}
	return s.repo.Expenses.Create(userID, categoryID, title, description, amount, expenseDate)
}

func (s *ExpenseServiceImpl) UpdateExpense(expenseID, requestingUserID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if categoryID <= 0 {
		return nil, fmt.Errorf("category_id is required")
	}
	if expenseDate.IsZero() {
		return nil, fmt.Errorf("expense_date is required")
	}

	existing, err := s.repo.Expenses.GetByID(expenseID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}
	if existing.UserID != requestingUserID {
		return nil, ErrForbidden
	}

	return s.repo.Expenses.Update(expenseID, title, description, amount, categoryID, expenseDate)
}

func (s *ExpenseServiceImpl) DeleteExpense(expenseID, requestingUserID int) error {
	existing, err := s.repo.Expenses.GetByID(expenseID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	if existing.UserID != requestingUserID {
		return ErrForbidden
	}
	return s.repo.Expenses.Delete(expenseID)
}

func (s *ExpenseServiceImpl) ListExpenses(userID, page, pageSize int) ([]models.Expense, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	return s.repo.Expenses.ListByUser(userID, pageSize, offset)
}
