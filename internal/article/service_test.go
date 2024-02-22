package article

import (
	"github.com/stretchr/testify/assert"
	"go-clean-architecture/internal/domain"
	"go-clean-architecture/mocks"
	"gorm.io/gorm"
	"testing"
)

func TestArticleService_Fetch(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)
	mocksArticleList := make([]*domain.Article, 0)
	mocksArticleList = append(mocksArticleList, &domain.Article{
		ID:      1,
		Title:   "Title 1",
		Content: "Content 1",
	})

	t.Run("success", func(t *testing.T) {
		mockArticleRepository.On("Fetch", uint(1), uint(10), &domain.Article{}).
			Return(mocksArticleList, uint(2), nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		articles, nextCursor, err := articleSvc.Fetch(uint(1), uint(10), &domain.Article{})
		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, uint(2), nextCursor)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("success-zero-size", func(t *testing.T) {
		mockArticleRepository.On("Fetch", uint(1), uint(10), &domain.Article{}).
			Return(mocksArticleList, uint(2), nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		articles, nextCursor, err := articleSvc.Fetch(uint(1), uint(10), &domain.Article{})
		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, uint(2), nextCursor)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepository.On("Fetch", uint(1), uint(10), &domain.Article{}).
			Return(nil, uint(0), assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		articles, nextCursor, err := articleSvc.Fetch(uint(1), uint(10), &domain.Article{})
		assert.Error(t, err)
		assert.Nil(t, articles)
		assert.Equal(t, uint(0), nextCursor)
	})
}

func TestArticleService_GetByID(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)
	mockArticle := &domain.Article{
		ID:      1,
		Title:   "Title 1",
		Content: "Content 1",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepository.On("GetByID", uint(1)).
			Return(mockArticle, nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		article, err := articleSvc.GetByID(uint(1))
		assert.NoError(t, err)
		assert.NotNil(t, article)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("error-not-found", func(t *testing.T) {
		mockArticleRepository.On("GetByID", uint(1)).
			Return(nil, gorm.ErrRecordNotFound).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		article, err := articleSvc.GetByID(uint(1))
		assert.Error(t, err)
		assert.Nil(t, article)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepository.On("GetByID", uint(1)).
			Return(nil, assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		article, err := articleSvc.GetByID(uint(1))
		assert.Error(t, err)
		assert.Nil(t, article)
	})
}

func TestArticleService_Count(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)

	t.Run("success", func(t *testing.T) {
		mockArticleRepository.On("Count", &domain.Article{}).
			Return(int64(10), nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		count, err := articleSvc.Count(&domain.Article{})
		assert.NoError(t, err)
		assert.Equal(t, int64(10), count)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepository.On("Count", &domain.Article{}).
			Return(int64(0), assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		count, err := articleSvc.Count(&domain.Article{})
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestArticleService_GetByTitle(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)
	mocksArticleList := make([]*domain.Article, 0)
	mocksArticleList = append(mocksArticleList, &domain.Article{
		ID:      1,
		Title:   "Title 1",
		Content: "Content 1",
	})

	t.Run("success", func(t *testing.T) {
		mockArticleRepository.On("GetByTitle", "Title").
			Return(mocksArticleList, nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		articles, err := articleSvc.GetByTitle("Title")
		assert.NoError(t, err)
		assert.NotNil(t, articles)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepository.On("GetByTitle", "Title").
			Return(nil, assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		articles, err := articleSvc.GetByTitle("Title")
		assert.Error(t, err)
		assert.Nil(t, articles)
	})
}

func TestArticleService_Store(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockArticle := &domain.Article{
		ID:       1,
		Title:    "Title 1",
		Content:  "Content 1",
		AuthorID: 1,
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", uint(1)).
			Return(&domain.Author{}, nil).Once()
		mockArticleRepository.On("Store", mockArticle).
			Return(nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, mockAuthorRepository)
		err := articleSvc.Store(mockArticle)
		assert.NoError(t, err)

		mockArticleRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("error-author-not-found", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", uint(1)).
			Return(nil, gorm.ErrRecordNotFound).Once()

		articleSvc := NewArticleService(mockArticleRepository, mockAuthorRepository)
		err := articleSvc.Store(mockArticle)
		assert.Error(t, err)

		mockArticleRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("error-author-failed", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", uint(1)).
			Return(nil, assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, mockAuthorRepository)
		err := articleSvc.Store(mockArticle)
		assert.Error(t, err)

		mockArticleRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", uint(1)).
			Return(&domain.Author{}, nil).Once()
		mockArticleRepository.On("Store", mockArticle).
			Return(assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, mockAuthorRepository)
		err := articleSvc.Store(mockArticle)
		assert.Error(t, err)

		mockArticleRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})
}

func TestArticleService_Update(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)
	mockArticle := &domain.Article{
		ID:      1,
		Title:   "Title 1",
		Content: "Content 1",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepository.On("Update", mockArticle).
			Return(nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		err := articleSvc.Update(mockArticle)
		assert.NoError(t, err)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepository.On("Update", mockArticle).
			Return(assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		err := articleSvc.Update(mockArticle)
		assert.Error(t, err)

		mockArticleRepository.AssertExpectations(t)
	})
}

func TestArticleService_Delete(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)

	t.Run("success", func(t *testing.T) {
		mockArticleRepository.On("Delete", uint(1)).
			Return(nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		err := articleSvc.Delete(uint(1))
		assert.NoError(t, err)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepository.On("Delete", uint(1)).
			Return(assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		err := articleSvc.Delete(uint(1))
		assert.Error(t, err)

		mockArticleRepository.AssertExpectations(t)
	})
}

func TestArticleService_GetByAuthorID(t *testing.T) {
	mockArticleRepository := new(mocks.ArticleRepository)
	mocksArticleList := make([]*domain.Article, 0)
	mocksArticleList = append(mocksArticleList, &domain.Article{
		ID:      1,
		Title:   "Title 1",
		Content: "Content 1",
	})

	t.Run("success", func(t *testing.T) {
		mockArticleRepository.On("GetByAuthorID", uint(1)).
			Return(mocksArticleList, nil).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		articles, err := articleSvc.GetByAuthorID(uint(1))
		assert.NoError(t, err)
		assert.NotNil(t, articles)

		mockArticleRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepository.On("GetByAuthorID", uint(1)).
			Return(nil, assert.AnError).Once()

		articleSvc := NewArticleService(mockArticleRepository, nil)
		articles, err := articleSvc.GetByAuthorID(uint(1))
		assert.Error(t, err)
		assert.Nil(t, articles)
	})
}
