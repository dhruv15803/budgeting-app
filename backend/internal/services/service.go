package services

import (
	"time"

	"github.com/dhruv15803/budgeting-app/internal/auth"
	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
	"github.com/dhruv15803/budgeting-app/internal/worker"
)

type Service struct {
	Users    UserService
	Expenses ExpenseService
}

func NewService(repo *repositories.Repository, cfg *config.Config, jwt *auth.JWTSigner, q *worker.Queue) *Service {
	return &Service{
		Users:    NewUserService(repo, cfg, jwt, q),
		Expenses: NewExpenseService(repo),
	}
}

type UserService interface {
	GetMe(userID int) (*models.User, error)
	DeleteUserById(id int) error
	Register(email string, password string, username *string) error
	Login(email string, password string) (token string, err error)
	VerifyEmail(rawToken string) (token string, err error)
}

type ExpenseService interface {
	CreateExpense(userID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error)
	UpdateExpense(expenseID, requestingUserID, categoryID int, title string, description *string, amount float64, expenseDate time.Time) (*models.Expense, error)
	DeleteExpense(expenseID, requestingUserID int) error
	ListExpenses(userID, page, pageSize int) ([]models.Expense, int, error)
}
