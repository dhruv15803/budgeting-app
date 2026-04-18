package services

import "github.com/dhruv15803/budgeting-app/internal/repositories"

type UserServiceImpl struct {
	repo *repositories.Repository
}

func NewUserService(repo *repositories.Repository) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (u *UserServiceImpl) DeleteUserById(id int) error {
	return u.repo.Users.DeleteUserById(id)
}
