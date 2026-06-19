package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAll(all bool) ([]models.Article, error)
	GetBySlug(slug string) (*models.Article, error)
	Create(article *models.Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) GetAll(all bool) ([]models.Article, error) {
	var articles []models.Article
	query := r.db.Preload("Author").Where("deleted_at IS NULL")
	if !all {
		query = query.Where("is_published = ?", true)
	}
	err := query.Order("published_at DESC, created_at DESC").Find(&articles).Error
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

func (r *articleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}
