package mocks

import (
	"github.com/stretchr/testify/mock"
	"go-clean-architecture/internal/domain"
)

type AuthorRepository struct {
	mock.Mock
}

func (m *AuthorRepository) GetByID(id uint) (*domain.Author, error) {
	ret := m.Called(id)

	var r0 *domain.Author
	if rf, ok := ret.Get(0).(func(uint) *domain.Author); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Author)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *AuthorRepository) Count() (int64, error) {
	ret := m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *AuthorRepository) Store(author *domain.Author) error {
	ret := m.Called(author)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Author) error); ok {
		r0 = rf(author)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
