package models

import (
	"time"

	"github.com/google/uuid"
)

type ArticleComment struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ArticleID uuid.UUID  `gorm:"type:uuid;not null;index" json:"articleId"`
	Article   *Article   `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	Email     string     `gorm:"type:varchar(255)" json:"email"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}
