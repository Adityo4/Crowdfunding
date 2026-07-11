package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ArticleCommentService interface {
	GetCommentsByArticleSlug(slug string) ([]models.ArticleComment, error)
	CreateComment(slug string, name, email, content string) (*models.ArticleComment, error)
}

type articleCommentService struct {
	commentRepo repository.ArticleCommentRepository
	articleRepo repository.ArticleRepository
}

func NewArticleCommentService(commentRepo repository.ArticleCommentRepository, articleRepo repository.ArticleRepository) ArticleCommentService {
	return &articleCommentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

func (s *articleCommentService) GetCommentsByArticleSlug(slug string) ([]models.ArticleComment, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}
	return s.commentRepo.GetByArticleID(article.ID)
}

func (s *articleCommentService) CreateComment(slug string, name, email, content string) (*models.ArticleComment, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	comment := &models.ArticleComment{
		ID:        uuid.New(),
		ArticleID: article.ID,
		Name:      name,
		Email:     email,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}
