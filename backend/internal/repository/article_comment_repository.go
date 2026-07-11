package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleCommentRepository interface {
	GetByArticleID(articleID uuid.UUID) ([]models.ArticleComment, error)
	Create(comment *models.ArticleComment) error
}

type articleCommentRepository struct {
	db *gorm.DB
}

func NewArticleCommentRepository(db *gorm.DB) ArticleCommentRepository {
	return &articleCommentRepository{db}
}

func (r *articleCommentRepository) GetByArticleID(articleID uuid.UUID) ([]models.ArticleComment, error) {
	var comments []models.ArticleComment
	err := r.db.Where("article_id = ? AND deleted_at IS NULL", articleID).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *articleCommentRepository) Create(comment *models.ArticleComment) error {
	return r.db.Create(comment).Error
}
