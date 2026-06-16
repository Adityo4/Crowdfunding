package service

import (
	"backend/internal/models"
	"backend/internal/repository"
)

type CategoryService interface {
	GetCategories() ([]models.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo}
}

func (s *categoryService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}
