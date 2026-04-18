package services

import "github.com/dhruv15803/budgeting-app/internal/repositories"

type Service struct {
	Users UserService
}

func NewService(repo *repositories.Repository) *Service {
	return &Service{
		Users: NewUserService(repo),
	}
}

type UserService interface {
	DeleteUserById(id int) error
}
