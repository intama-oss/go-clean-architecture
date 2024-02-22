package author

import (
	"go-clean-architecture/internal/domain"
	"gorm.io/gorm"
)

type mysqlAuthorRepository struct {
	db *gorm.DB
}

func NewMysqlAuthorRepository(db *gorm.DB) domain.AuthorRepository {
	return &mysqlAuthorRepository{db: db}
}

func (r *mysqlAuthorRepository) GetByID(id uint) (*domain.Author, error) {
	var author domain.Author
	if err := r.db.First(&author, id).Error; err != nil {
		return nil, err
	}
	return &author, nil
}

func (r *mysqlAuthorRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Author{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *mysqlAuthorRepository) Store(author *domain.Author) error {
	return r.db.Create(author).Error
}
