package models

import (
	"time"

	"github.com/google/uuid"
)

type Donation struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CharityID        uuid.UUID  `gorm:"type:uuid;not null" json:"charityId"`
	Charity          *Charity   `gorm:"foreignKey:CharityID" json:"charity,omitempty"`
	UserID           *uuid.UUID `gorm:"type:uuid" json:"userId,omitempty"`
	User             *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Amount           float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	Status           string     `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	PaymentMethod    string     `gorm:"type:varchar(100)" json:"paymentMethod"`
	PaymentReference string     `gorm:"type:varchar(255)" json:"paymentReference,omitempty"`
	Anonymous        bool       `gorm:"type:boolean;default:false;not null" json:"anonymous"`
	Message          string     `gorm:"type:text" json:"message"`
	SnapToken        string     `gorm:"-" json:"snapToken,omitempty"`
	RedirectURL      string     `gorm:"-" json:"redirectUrl,omitempty"`
	CreatedAt        time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
}

type CreateDonationRequest struct {
	CharityID     uuid.UUID `json:"charityId" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=1000"`
	PaymentMethod string    `json:"paymentMethod" binding:"required"`
	Anonymous     bool      `json:"anonymous"`
	Message       string    `json:"message"`
}

type CallbackRequest struct {
	PaymentReference string `json:"paymentReference" binding:"required"`
	DonationID       string `json:"donationId" binding:"required,uuid"`
	Status           string `json:"status" binding:"required"`
}
