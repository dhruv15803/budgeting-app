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
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:                db,
		Users:             NewUserRepo(db),
		EmailVerification: NewEmailVerificationRepo(db),
		Expenses:          NewExpenseRepo(db),
	}
}

func (r *Repository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

type UserRepository interface {
	GetByID(id int) (*models.User, error)
	DeleteUserById(id int) error
	GetByEmail(email string) (*models.User, error)
	GetByEmailTx(tx *sqlx.Tx, email string) (*models.User, error)
	CreateUserTx(tx *sqlx.Tx, email string, username *string, passwordHash string) (int, error)
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
}
