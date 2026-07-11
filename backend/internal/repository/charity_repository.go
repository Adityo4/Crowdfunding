package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CharityRepository interface {
	Create(charity *models.Charity) error
	CreateImage(img *models.CharityImage) error
	GetAll(q string, categorySlug string, status string, sort string, offset, limit int) ([]models.Charity, int64, error)
	GetBySlug(slug string) (*models.Charity, error)
	GetByID(id uuid.UUID) (*models.Charity, error)
	Update(charity *models.Charity) error
}

type charityRepository struct {
	db *gorm.DB
}

func NewCharityRepository(db *gorm.DB) CharityRepository {
	return &charityRepository{db}
}

func (r *charityRepository) Create(charity *models.Charity) error {
	return r.db.Create(charity).Error
}

func (r *charityRepository) CreateImage(img *models.CharityImage) error {
	return r.db.Create(img).Error
}

func (r *charityRepository) GetAll(q string, categorySlug string, status string, sort string, offset, limit int) ([]models.Charity, int64, error) {
	var charities []models.Charity
	var total int64

	query := r.db.Model(&models.Charity{}).Preload("Category")

	// Filter Search Query
	if q != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	// Filter Category Slug
	if categorySlug != "" {
		query = query.Joins("JOIN categories ON categories.id = charities.category_id").
			Where("categories.slug = ?", categorySlug)
	}

	// Filter Status
	if status != "" {
		if status != "all" {
			query = query.Where("charities.status = ?", status)
		}
	} else {
		query = query.Where("charities.status = ?", "active") // default active
	}

	query = query.Where("charities.deleted_at IS NULL")

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	switch sort {
	case "target_amount":
		query = query.Order("target_amount DESC")
	case "urgent":
		query = query.Order("end_date ASC")
	case "newest":
		fallthrough
	default:
		query = query.Order("created_at DESC")
	}

	// Pagination
	err := query.Offset(offset).Limit(limit).Find(&charities).Error
	if err != nil {
		return nil, 0, err
	}

	return charities, total, nil
}

func (r *charityRepository) GetBySlug(slug string) (*models.Charity, error) {
	var charity models.Charity
	err := r.db.Preload("Category").Preload("User").Preload("Images").
		Preload("Donations", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Where("status = ?", "paid").Order("created_at DESC")
		}).
		Where("slug = ? AND deleted_at IS NULL", slug).First(&charity).Error
	if err != nil {
		return nil, err
	}
	return &charity, nil
}

func (r *charityRepository) GetByID(id uuid.UUID) (*models.Charity, error) {
	var charity models.Charity
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&charity).Error
	if err != nil {
		return nil, err
	}
	return &charity, nil
}

func (r *charityRepository) Update(charity *models.Charity) error {
	return r.db.Save(charity).Error
}
