package repositories

import (
	"time"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db                *sqlx.DB
	Users             UserRepository
	EmailVerification EmailVerificationRepository
	Expenses          ExpenseRepository
	RecurringExpenses RecurringExpenseRepository
	Budgets           BudgetRepository
	ExpenseCategories ExpenseCategoryRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:                db,
		Users:             NewUserRepo(db),
		EmailVerification: NewEmailVerificationRepo(db),
		Expenses:          NewExpenseRepo(db),
		RecurringExpenses: NewRecurringExpenseRepo(db),
		Budgets:           NewBudgetRepo(db),
		ExpenseCategories: NewExpenseCategoryRepo(db),
	}
}

func (r *Repository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

type UserRepository interface {
	GetByID(id int) (*models.User, error)
	DeleteUserById(id int) error
	GetByGoogleSub(googleSub string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByEmailTx(tx *sqlx.Tx, email string) (*models.User, error)
	CreateUserTx(tx *sqlx.Tx, email string, username *string, passwordHash string) (int, error)
	CreateGoogleUserTx(tx *sqlx.Tx, email string, googleSub string, imageURL *string) (int, error)
	LinkGoogleIdentityTx(tx *sqlx.Tx, userID int, googleSub string, imageURL *string) error
	UpdateUnverifiedCredentialsTx(tx *sqlx.Tx, userID int, username *string, passwordHash string) error
}

type EmailVerificationRepository interface {
	DeleteByUserIDTx(tx *sqlx.Tx, userID int) error
	InsertTx(tx *sqlx.Tx, userID int, tokenHash string, expiresAt time.Time) error
	VerifyTokenTx(tx *sqlx.Tx, tokenHash string) (*models.User, error)
}

type ExpenseRepository interface {
	Create(userID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error)
	GetByID(id int) (*models.Expense, error)
	Update(id int, title string, description *string, amount float64, categoryID int, expenseDate time.Time) (*models.Expense, error)
	Delete(id int) error
	ListByUser(userID int, f models.ExpenseFilter) ([]models.Expense, int, error)
	CreateFromRecurringTx(tx *sqlx.Tx, r *models.RecurringExpense, date time.Time) (bool, error)
}

type RecurringExpenseRepository interface {
	Create(userID, categoryID int, title string, description *string, amount float64, startDate time.Time, endDate *time.Time, frequency string) (*models.RecurringExpense, error)
	GetByID(id int) (*models.RecurringExpense, error)
	Update(id int, title string, description *string, amount float64, categoryID int, startDate time.Time, endDate *time.Time, frequency string, isActive bool, nextOccurrence time.Time) (*models.RecurringExpense, error)
	Delete(id int) error
	ListByUser(userID int, f models.RecurringExpenseFilter) ([]models.RecurringExpense, int, error)
	ListDue(today time.Time) ([]models.RecurringExpense, error)
	AdvanceTx(tx *sqlx.Tx, id int, newNext time.Time, deactivate bool) error
}

type ExpenseCategoryRepository interface {
	ListAll() ([]models.ExpenseCategory, error)
}

type BudgetRepository interface {
	// monthly_budgets
	CreateBudget(userID int, budgetMonth time.Time, totalAmount float64) (*models.MonthlyBudget, error)
	GetBudgetByMonth(userID int, budgetMonth time.Time) (*models.MonthlyBudget, error)
	GetBudgetByID(id int) (*models.MonthlyBudget, error)
	UpdateBudgetTotal(id int, totalAmount float64) (*models.MonthlyBudget, error)
	DeleteBudget(id int) error
	ListBudgets(userID int, f models.BudgetFilter) ([]models.MonthlyBudget, int, error)

	// monthly_category_budgets
	UpsertCategoryBudget(budgetID, categoryID int, amount float64) (*models.MonthlyBudgetCategory, error)
	DeleteCategoryBudget(budgetID, categoryID int) error
	ListCategoryBudgets(budgetID int) ([]models.MonthlyBudgetCategory, error)
	SumCategoryAllocations(budgetID int) (float64, error)
	BulkSetCategoryBudgetsTx(tx *sqlx.Tx, budgetID int, items []models.BulkCategoryItem) error

	// overview
	GetCategorySpending(userID, budgetID int, from, to time.Time) ([]models.CategoryBudgetWithSpending, error)
	GetTotalSpending(userID int, from, to time.Time) (float64, error)
}
