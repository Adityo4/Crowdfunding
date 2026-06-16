package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DonationRepository interface {
	Create(donation *models.Donation) error
	GetByID(id uuid.UUID) (*models.Donation, error)
	Update(donation *models.Donation) error
}

type donationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) DonationRepository {
	return &donationRepository{db}
}

func (r *donationRepository) Create(donation *models.Donation) error {
	return r.db.Create(donation).Error
}

func (r *donationRepository) GetByID(id uuid.UUID) (*models.Donation, error) {
	var donation models.Donation
	err := r.db.Preload("Charity").Where("id = ?", id).First(&donation).Error
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

func (r *donationRepository) Update(donation *models.Donation) error {
	return r.db.Save(donation).Error
}
