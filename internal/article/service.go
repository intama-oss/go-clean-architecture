package article

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-clean-architecture/internal/domain"
	"gorm.io/gorm"
)

type articleService struct {
	articleRepo domain.ArticleRepository
	authorRepo  domain.AuthorRepository
}

func NewArticleService(article domain.ArticleRepository, author domain.AuthorRepository) domain.ArticleService {
	return &articleService{
		articleRepo: article,
		authorRepo:  author,
	}
}

func (a *articleService) Fetch(page uint, size uint, filter *domain.Article) ([]*domain.Article, uint, error) {
	articles, nextCursor, err := a.articleRepo.Fetch(page, size, filter)
	if err != nil {
		return nil, 0, err
	}

	return articles, nextCursor, nil
}

func (a *articleService) GetByID(id uint) (*domain.Article, error) {
	article, err := a.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	return article, nil
}

func (a *articleService) Count(filter *domain.Article) (int64, error) {
	count, err := a.articleRepo.Count(filter)
	return count, err
}

func (a *articleService) GetByTitle(title string) ([]*domain.Article, error) {
	articles, err := a.articleRepo.GetByTitle(title)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (a *articleService) Store(article *domain.Article) error {
	author, err := a.authorRepo.GetByID(article.AuthorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return err
	}

	article.Author = author
	return a.articleRepo.Store(article)
}

func (a *articleService) Update(article *domain.Article) error {
	return a.articleRepo.Update(article)
}

func (a *articleService) Delete(id uint) error {
	return a.articleRepo.Delete(id)
}

func (a *articleService) GetByAuthorID(authorID uint) ([]*domain.Article, error) {
	articles, err := a.articleRepo.GetByAuthorID(authorID)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
