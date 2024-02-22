package domain

import (
	"time"
)

type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title" gorm:"type:varchar(255)"`
	Content   string    `json:"content" gorm:"type:text"`
	AuthorID  uint      `json:"authorId" gorm:"index"`
	Author    *Author   `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type ArticleStoreRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	AuthorID uint   `json:"authorId" validate:"required"`
}

type ArticleUpdateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleRepository interface {
	Fetch(page uint, size uint, filter *Article) ([]*Article, uint, error)
	GetByID(id uint) (*Article, error)
	Count(filter *Article) (int64, error)
	GetByAuthorID(authorID uint) ([]*Article, error)
	GetByTitle(title string) ([]*Article, error)
	Store(article *Article) error
	Update(article *Article) error
	Delete(id uint) error
}

type ArticleService interface {
	Fetch(page uint, size uint, filter *Article) ([]*Article, uint, error)
	GetByID(id uint) (*Article, error)
	Count(filter *Article) (int64, error)
	GetByTitle(title string) ([]*Article, error)
	GetByAuthorID(authorID uint) ([]*Article, error)
	Store(article *Article) error
	Update(article *Article) error
	Delete(id uint) error
}
