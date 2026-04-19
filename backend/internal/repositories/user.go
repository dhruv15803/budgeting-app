package repositories

import (
	"database/sql"
	"errors"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) DeleteUserById(id int) error {
	_, err := u.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetByEmail(email string) (*models.User, error) {
	var out models.User
	err := u.db.Get(&out, `
		SELECT id, email, username, password, image_url, role::text AS role,
		       is_verified, verified_at, created_at, updated_at
		FROM users WHERE email = $1
	`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (u *UserRepo) GetByEmailTx(tx *sqlx.Tx, email string) (*models.User, error) {
	var out models.User
	err := tx.Get(&out, `
		SELECT id, email, username, password, image_url, role::text AS role,
		       is_verified, verified_at, created_at, updated_at
		FROM users WHERE email = $1
	`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (u *UserRepo) CreateUserTx(tx *sqlx.Tx, email string, username *string, passwordHash string) (int, error) {
	var id int
	err := tx.QueryRow(`
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`, email, username, passwordHash).Scan(&id)
	return id, err
}

func (u *UserRepo) UpdateUnverifiedCredentialsTx(tx *sqlx.Tx, userID int, username *string, passwordHash string) error {
	res, err := tx.Exec(`
		UPDATE users
		SET username = COALESCE($2, username),
		    password = $3,
		    updated_at = NOW()
		WHERE id = $1 AND is_verified = false
	`, userID, username, passwordHash)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("user not found or already verified")
	}
	return nil
}
