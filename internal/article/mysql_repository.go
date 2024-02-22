package article

import (
	"go-clean-architecture/internal/domain"
	"gorm.io/gorm"
)

type mysqlArticleRepository struct {
	db *gorm.DB
}

func NewMysqlArticleRepository(db *gorm.DB) domain.ArticleRepository {
	return &mysqlArticleRepository{db: db}
}

func (r *mysqlArticleRepository) Fetch(page uint, size uint, filter *domain.Article) ([]*domain.Article, uint, error) {
	var articles []*domain.Article

	offset := (page - 1) * size
	query := r.db

	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}

	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(size)).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor uint
	if len(articles) > 0 {
		nextCursor = page + 1 // next page
	}

	return articles, nextCursor, nil
}

func (r *mysqlArticleRepository) GetByID(id uint) (*domain.Article, error) {
	var article *domain.Article
	if err := r.db.Preload("Author").First(&article, id).Error; err != nil {
		return nil, err
	}
	return article, nil
}

func (r *mysqlArticleRepository) Count(filter *domain.Article) (int64, error) {
	var count int64
	query := r.db.Model(&domain.Article{})

	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *mysqlArticleRepository) Store(article *domain.Article) error {
	return r.db.Create(article).Error
}

func (r *mysqlArticleRepository) Update(article *domain.Article) error {
	return r.db.Updates(article).Error
}

func (r *mysqlArticleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, id).Error
}

func (r *mysqlArticleRepository) GetByAuthorID(authorID uint) ([]*domain.Article, error) {
	var articles []*domain.Article
	if err := r.db.Where("author_id = ?", authorID).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *mysqlArticleRepository) GetByTitle(title string) ([]*domain.Article, error) {
	var articles []*domain.Article
	if err := r.db.Where("title LIKE ?", "%"+title+"%").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}
