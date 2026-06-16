package service

import (
	"backend/internal/models"
	"backend/internal/repository"
)

type ArticleService interface {
	GetArticles() ([]models.Article, error)
	GetArticleBySlug(slug string) (*models.Article, error)
}

type articleService struct {
	articleRepo repository.ArticleRepository
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{articleRepo}
}

func (s *articleService) GetArticles() ([]models.Article, error) {
	return s.articleRepo.GetAll()
}

func (s *articleService) GetArticleBySlug(slug string) (*models.Article, error) {
	return s.articleRepo.GetBySlug(slug)
}
