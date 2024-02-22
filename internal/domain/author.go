package domain

import "time"

type Author struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255)" validate:"required"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type AuthorRepository interface {
	GetByID(id uint) (*Author, error)
	Count() (int64, error)
	Store(author *Author) error
}
