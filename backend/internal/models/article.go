package models

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AuthorID      uuid.UUID  `gorm:"type:uuid;not null" json:"authorId"`
	Author        User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Title         string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug          string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content       string     `gorm:"type:text;not null" json:"content"`
	CoverImageUrl string     `gorm:"type:varchar(255)" json:"coverImageUrl"`
	IsPublished   bool       `gorm:"type:boolean;default:false;not null" json:"isPublished"`
	PublishedAt   *time.Time `json:"publishedAt,omitempty"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"index" json:"-"`
}
