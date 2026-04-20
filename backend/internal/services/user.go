package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/auth"
	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/dhruv15803/budgeting-app/internal/models"
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

func (u *UserServiceImpl) GetMe(userID int) (*models.User, error) {
	user, err := u.repo.Users.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	return user, nil
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
	if !user.Password.Valid || user.Password.String == "" {
		return "", ErrInvalidCredentials
	}
	if err := auth.ComparePassword(user.Password.String, password); err != nil {
		return "", ErrInvalidCredentials
	}
	if !user.IsVerified {
		return "", ErrNotVerified
	}
	return u.jwt.SignAccessToken(user.ID, user.Email, user.Role)
}

func (u *UserServiceImpl) LoginWithGoogle(ctx context.Context, credential string) (string, error) {
	if strings.TrimSpace(u.cfg.GoogleOAuthClientID) == "" {
		return "", ErrGoogleAuthDisabled
	}
	claims, err := auth.ParseGoogleIDToken(ctx, strings.TrimSpace(credential), u.cfg.GoogleOAuthClientID)
	if err != nil {
		return "", ErrInvalidGoogleToken
	}
	if !claims.EmailVerified {
		return "", ErrGoogleEmailNotVerified
	}
	email := normalizeEmail(claims.Email)
	sub := claims.Sub
	var picture *string
	if claims.Picture != "" {
		picture = &claims.Picture
	}

	bySub, err := u.repo.Users.GetByGoogleSub(sub)
	if err != nil {
		return "", err
	}
	if bySub != nil {
		return u.jwt.SignAccessToken(bySub.ID, bySub.Email, bySub.Role)
	}

	byEmail, err := u.repo.Users.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if byEmail != nil {
		if byEmail.GoogleSub.Valid && byEmail.GoogleSub.String != "" && byEmail.GoogleSub.String != sub {
			return "", ErrGoogleAccountConflict
		}
		if byEmail.GoogleSub.Valid && byEmail.GoogleSub.String == sub {
			return u.jwt.SignAccessToken(byEmail.ID, byEmail.Email, byEmail.Role)
		}
	}

	tx, err := u.repo.BeginTx()
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	if byEmail == nil {
		id, err := u.repo.Users.CreateGoogleUserTx(tx, email, sub, picture)
		if err != nil {
			return "", err
		}
		if err := tx.Commit(); err != nil {
			return "", err
		}
		return u.jwt.SignAccessToken(id, email, "user")
	}

	if err := u.repo.Users.LinkGoogleIdentityTx(tx, byEmail.ID, sub, picture); err != nil {
		if errors.Is(err, repositories.ErrGoogleLinkConflict) {
			return "", ErrGoogleAccountConflict
		}
		return "", err
	}
	if err := u.repo.EmailVerification.DeleteByUserIDTx(tx, byEmail.ID); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return u.jwt.SignAccessToken(byEmail.ID, byEmail.Email, byEmail.Role)
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
