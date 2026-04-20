package repositories

import (
	"database/sql"
	"errors"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/jmoiron/sqlx"
)

// ErrGoogleLinkConflict is returned when a user row cannot be linked to the given Google sub (e.g. different sub already stored).
var ErrGoogleLinkConflict = errors.New("google link conflict")

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) GetByID(id int) (*models.User, error) {
	var out models.User
	err := u.db.Get(&out, `
		SELECT id, email, username, password, google_sub, image_url, role::text AS role,
		       is_verified, verified_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (u *UserRepo) DeleteUserById(id int) error {
	_, err := u.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetByGoogleSub(googleSub string) (*models.User, error) {
	var out models.User
	err := u.db.Get(&out, `
		SELECT id, email, username, password, google_sub, image_url, role::text AS role,
		       is_verified, verified_at, created_at, updated_at
		FROM users WHERE google_sub = $1
	`, googleSub)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &out, err
}

func (u *UserRepo) GetByEmail(email string) (*models.User, error) {
	var out models.User
	err := u.db.Get(&out, `
		SELECT id, email, username, password, google_sub, image_url, role::text AS role,
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
		SELECT id, email, username, password, google_sub, image_url, role::text AS role,
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

func (u *UserRepo) CreateGoogleUserTx(tx *sqlx.Tx, email string, googleSub string, imageURL *string) (int, error) {
	var id int
	err := tx.QueryRow(`
		INSERT INTO users (email, username, password, google_sub, image_url, is_verified, verified_at)
		VALUES ($1, NULL, NULL, $2, $3, true, NOW())
		RETURNING id
	`, email, googleSub, imageURL).Scan(&id)
	return id, err
}

func (u *UserRepo) LinkGoogleIdentityTx(tx *sqlx.Tx, userID int, googleSub string, imageURL *string) error {
	res, err := tx.Exec(`
		UPDATE users SET
			google_sub = $2,
			image_url = COALESCE($3, image_url),
			is_verified = true,
			verified_at = COALESCE(verified_at, NOW()),
			updated_at = NOW()
		WHERE id = $1 AND (google_sub IS NULL OR google_sub = $2)
	`, userID, googleSub, imageURL)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrGoogleLinkConflict
	}
	return nil
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
