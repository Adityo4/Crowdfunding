package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAll() ([]models.Article, error)
	GetBySlug(slug string) (*models.Article, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) GetAll() ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Preload("Author").Where("is_published = ? AND deleted_at IS NULL", true).
		Order("published_at DESC").Find(&articles).Error
	return articles, err
}

func (r *articleRepository) GetBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Author").Where("slug = ? AND is_published = ? AND deleted_at IS NULL", slug, true).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}
