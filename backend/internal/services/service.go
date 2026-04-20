package services

import (
	"context"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/auth"
	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
	"github.com/dhruv15803/budgeting-app/internal/worker"
)

type Service struct {
	Users              UserService
	Expenses           ExpenseService
	RecurringExpenses  RecurringExpenseService
	Budgets            BudgetService
	ExpenseCategories  ExpenseCategoryService
}

func NewService(repo *repositories.Repository, cfg *config.Config, jwt *auth.JWTSigner, q *worker.Queue, logger func(format string, args ...interface{})) *Service {
	return &Service{
		Users:             NewUserService(repo, cfg, jwt, q),
		Expenses:          NewExpenseService(repo),
		RecurringExpenses: NewRecurringExpenseService(repo, logger),
		Budgets:           NewBudgetService(repo),
		ExpenseCategories: NewExpenseCategoryService(repo),
	}
}

type UserService interface {
	GetMe(userID int) (*models.User, error)
	DeleteUserById(id int) error
	Register(email string, password string, username *string) error
	Login(email string, password string) (token string, err error)
	LoginWithGoogle(ctx context.Context, credential string) (token string, err error)
	VerifyEmail(rawToken string) (token string, err error)
}

type ExpenseService interface {
	CreateExpense(userID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error)
	UpdateExpense(expenseID, requestingUserID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error)
	DeleteExpense(expenseID, requestingUserID int) error
	ListExpenses(userID int, f models.ExpenseFilter) ([]models.Expense, int, error)
}

type RecurringExpenseService interface {
	CreateRecurringExpense(userID, categoryID int, title string, description *string, amount float64, startDate time.Time, endDate *time.Time, frequency string) (*models.RecurringExpense, error)
	GetRecurringExpense(id, requestingUserID int) (*models.RecurringExpense, error)
	UpdateRecurringExpense(id, requestingUserID, categoryID int, title string, description *string, amount float64, startDate time.Time, endDate *time.Time, frequency string, isActive bool) (*models.RecurringExpense, error)
	DeleteRecurringExpense(id, requestingUserID int) error
	ListRecurringExpenses(userID int, f models.RecurringExpenseFilter) ([]models.RecurringExpense, int, error)
	RunDueGenerator(ctx context.Context) (int, error)
}

type ExpenseCategoryService interface {
	ListAll() ([]models.ExpenseCategory, error)
}

type BudgetService interface {
	CreateBudget(userID int, budgetMonth time.Time, totalAmount float64) (*models.MonthlyBudget, error)
	GetBudget(userID int, budgetMonth time.Time) (*BudgetOverviewResult, error)
	UpdateBudgetTotal(userID int, budgetMonth time.Time, totalAmount float64) (*models.MonthlyBudget, error)
	DeleteBudget(userID int, budgetMonth time.Time) error
	ListBudgets(userID int, f models.BudgetFilter) ([]models.MonthlyBudget, int, error)
	UpsertCategoryBudget(userID int, budgetMonth time.Time, categoryID int, amount float64) (*models.MonthlyBudgetCategory, error)
	DeleteCategoryBudget(userID int, budgetMonth time.Time, categoryID int) error
	BulkSetCategoryBudgets(userID int, budgetMonth time.Time, items []models.BulkCategoryItem) error
}
