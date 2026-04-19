package services

import (
	"github.com/dhruv15803/budgeting-app/internal/auth"
	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
	"github.com/dhruv15803/budgeting-app/internal/worker"
)

type Service struct {
	Users UserService
}

func NewService(repo *repositories.Repository, cfg *config.Config, jwt *auth.JWTSigner, q *worker.Queue) *Service {
	return &Service{
		Users: NewUserService(repo, cfg, jwt, q),
	}
}

type UserService interface {
	DeleteUserById(id int) error
	Register(email string, password string, username *string) error
	Login(email string, password string) (token string, err error)
	VerifyEmail(rawToken string) (token string, err error)
}
