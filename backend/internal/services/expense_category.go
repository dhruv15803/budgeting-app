package services

import (
	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
)

type expenseCategoryService struct {
	repo *repositories.Repository
}

func NewExpenseCategoryService(repo *repositories.Repository) ExpenseCategoryService {
	return &expenseCategoryService{repo: repo}
}

func (s *expenseCategoryService) ListAll() ([]models.ExpenseCategory, error) {
	return s.repo.ExpenseCategories.ListAll()
}
