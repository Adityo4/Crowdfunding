package models

import (
	"time"

	"github.com/google/uuid"
)

type Charity struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
	User             User           `gorm:"foreignKey:UserID" json:"creator,omitempty"`
	CategoryID       uuid.UUID      `gorm:"type:uuid;not null" json:"categoryId"`
	Category         Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Title            string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug             string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description      string         `gorm:"type:text;not null" json:"description"`
	TargetAmount     float64        `gorm:"type:decimal(15,2);not null" json:"targetAmount"`
	CurrentAmount    float64        `gorm:"type:decimal(15,2);default:0;not null" json:"currentAmount"`
	StartDate        time.Time      `gorm:"type:date;not null" json:"startDate"`
	EndDate          time.Time      `gorm:"type:date;not null" json:"endDate"`
	CoverImageUrl    string         `gorm:"type:varchar(255)" json:"coverImageUrl"`
	ContactPerson    string         `gorm:"type:varchar(255)" json:"contactPerson"`
	ContactEmail     string         `gorm:"type:varchar(255)" json:"contactEmail"`
	ContactPhone     string         `gorm:"type:varchar(50)" json:"contactPhone"`
	OrganizationName string         `gorm:"type:varchar(255)" json:"organizationName"`
	Status           string         `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	IsVerified       bool           `gorm:"type:boolean;default:false;not null" json:"isVerified"`
	Images           []CharityImage `gorm:"foreignKey:CharityID" json:"additionalImages,omitempty"`
	CreatedAt        time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt        *time.Time     `gorm:"index" json:"-"`
}

type CharityImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CharityID uuid.UUID `gorm:"type:uuid;not null" json:"charityId"`
	ImageUrl  string    `gorm:"type:varchar(255);not null" json:"imageUrl"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
}

type CreateCharityRequest struct {
	Title            string    `form:"title" binding:"required"`
	Description      string    `form:"description" binding:"required"`
	CategoryID       string    `form:"categoryId" binding:"required,uuid"`
	TargetAmount     float64   `form:"targetAmount" binding:"required,gt=10000"`
	StartDate        time.Time `form:"startDate" time_format:"2006-01-02" binding:"required"`
	EndDate          time.Time `form:"endDate" time_format:"2006-01-02" binding:"required"`
	ContactPerson    string    `form:"contactPerson" binding:"required"`
	ContactEmail     string    `form:"contactEmail" binding:"required,email"`
	ContactPhone     string    `form:"contactPhone"`
	OrganizationName string    `form:"organization"`
}
