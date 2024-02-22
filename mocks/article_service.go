package mocks

import (
	"github.com/stretchr/testify/mock"
	"go-clean-architecture/internal/domain"
)

type ArticleService struct {
	mock.Mock
}

func (m *ArticleService) Fetch(page uint, size uint, filter *domain.Article) ([]*domain.Article, uint, error) {
	ret := m.Called(page, size, filter)

	var r0 []*domain.Article
	if rf, ok := ret.Get(0).(func(page uint, size uint, filter *domain.Article) []*domain.Article); ok {
		r0 = rf(page, size, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Article)
		}
	}

	var r1 uint
	if rf, ok := ret.Get(1).(func(page uint, size uint, filter *domain.Article) uint); ok {
		r1 = rf(page, size, filter)
	} else {
		r1 = ret.Get(1).(uint)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(page uint, size uint, filter *domain.Article) error); ok {
		r2 = rf(page, size, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

func (m *ArticleService) GetByID(id uint) (*domain.Article, error) {
	ret := m.Called(id)

	var r0 *domain.Article
	if rf, ok := ret.Get(0).(func(id uint) *domain.Article); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(id uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ArticleService) Count(filter *domain.Article) (int64, error) {
	ret := m.Called(filter)

	var r0 int64
	if rf, ok := ret.Get(0).(func(filter *domain.Article) int64); ok {
		r0 = rf(filter)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(filter *domain.Article) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ArticleService) GetByTitle(title string) ([]*domain.Article, error) {
	ret := m.Called(title)

	var r0 []*domain.Article
	if rf, ok := ret.Get(0).(func(title string) []*domain.Article); ok {
		r0 = rf(title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(title string) error); ok {
		r1 = rf(title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ArticleService) GetByAuthorID(authorID uint) ([]*domain.Article, error) {
	ret := m.Called(authorID)

	var r0 []*domain.Article
	if rf, ok := ret.Get(0).(func(authorID uint) []*domain.Article); ok {
		r0 = rf(authorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(authorID uint) error); ok {
		r1 = rf(authorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ArticleService) Store(article *domain.Article) error {
	ret := m.Called(article)

	var r0 error
	if rf, ok := ret.Get(0).(func(article *domain.Article) error); ok {
		r0 = rf(article)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *ArticleService) Update(article *domain.Article) error {
	ret := m.Called(article)

	var r0 error
	if rf, ok := ret.Get(0).(func(article *domain.Article) error); ok {
		r0 = rf(article)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *ArticleService) Delete(id uint) error {
	ret := m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(id uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
