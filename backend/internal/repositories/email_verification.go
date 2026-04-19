package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type EmailVerificationRepo struct {
	db *sqlx.DB
}

func NewEmailVerificationRepo(db *sqlx.DB) *EmailVerificationRepo {
	return &EmailVerificationRepo{db: db}
}

func (e *EmailVerificationRepo) DeleteByUserIDTx(tx *sqlx.Tx, userID int) error {
	_, err := tx.Exec(`DELETE FROM email_verifications WHERE user_id = $1`, userID)
	return err
}

func (e *EmailVerificationRepo) InsertTx(tx *sqlx.Tx, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := tx.Exec(`
		INSERT INTO email_verifications (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt.UTC())
	return err
}

func (e *EmailVerificationRepo) VerifyTokenTx(tx *sqlx.Tx, tokenHash string) (*models.User, error) {
	var userID int
	err := tx.Get(&userID, `
		SELECT user_id FROM email_verifications
		WHERE token_hash = $1 AND expires_at > NOW()
	`, tokenHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		UPDATE users
		SET is_verified = true,
		    verified_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`DELETE FROM email_verifications WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	var u models.User
	err = tx.Get(&u, `
		SELECT id, email, username, password, image_url, role::text AS role,
		       is_verified, verified_at, created_at, updated_at
		FROM users WHERE id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
