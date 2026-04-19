package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/auth"
	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
	"github.com/dhruv15803/budgeting-app/internal/worker"
)

type UserServiceImpl struct {
	repo *repositories.Repository
	cfg  *config.Config
	jwt  *auth.JWTSigner
	q    *worker.Queue
}

func NewUserService(repo *repositories.Repository, cfg *config.Config, jwt *auth.JWTSigner, q *worker.Queue) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
		cfg:  cfg,
		jwt:  jwt,
		q:    q,
	}
}

func (u *UserServiceImpl) DeleteUserById(id int) error {
	return u.repo.Users.DeleteUserById(id)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validationPassword(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

func randomVerificationRawToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func verificationLink(baseURL string, rawToken string) string {
	base := strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
	return base + "/api/auth/verify-email?token=" + url.QueryEscape(rawToken)
}

func (u *UserServiceImpl) Register(email string, password string, username *string) error {
	email = normalizeEmail(email)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if err := validationPassword(password); err != nil {
		return err
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	tx, err := u.repo.BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	existing, err := u.repo.Users.GetByEmailTx(tx, email)
	if err != nil {
		return err
	}

	var userID int
	switch {
	case existing == nil:
		id, err := u.repo.Users.CreateUserTx(tx, email, username, passwordHash)
		if err != nil {
			return err
		}
		userID = id
	case existing.IsVerified:
		return ErrEmailTaken
	default:
		if err := u.repo.Users.UpdateUnverifiedCredentialsTx(tx, existing.ID, username, passwordHash); err != nil {
			return err
		}
		userID = existing.ID
	}

	if err := u.repo.EmailVerification.DeleteByUserIDTx(tx, userID); err != nil {
		return err
	}

	rawToken, err := randomVerificationRawToken()
	if err != nil {
		return err
	}
	tokenHash := auth.HashVerificationToken(rawToken)
	expiresAt := time.Now().UTC().Add(u.cfg.EmailVerificationTokenTTL)
	if err := u.repo.EmailVerification.InsertTx(tx, userID, tokenHash, expiresAt); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	job := worker.VerificationJob{
		ToEmail:         email,
		VerificationURL: verificationLink(u.cfg.EmailVerificationBaseURL, rawToken),
	}
	if !u.q.Submit(job) {
		log.Printf("verification email queue full; email not queued for %s", email)
	}
	return nil
}

func (u *UserServiceImpl) Login(email string, password string) (string, error) {
	email = normalizeEmail(email)
	user, err := u.repo.Users.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}
	if err := auth.ComparePassword(user.Password, password); err != nil {
		return "", ErrInvalidCredentials
	}
	if !user.IsVerified {
		return "", ErrNotVerified
	}
	return u.jwt.SignAccessToken(user.ID, user.Email, user.Role)
}

func (u *UserServiceImpl) VerifyEmail(rawToken string) (string, error) {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return "", ErrInvalidOrExpiredToken
	}
	tokenHash := auth.HashVerificationToken(rawToken)

	tx, err := u.repo.BeginTx()
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	user, err := u.repo.EmailVerification.VerifyTokenTx(tx, tokenHash)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidOrExpiredToken
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return u.jwt.SignAccessToken(user.ID, user.Email, user.Role)
}
