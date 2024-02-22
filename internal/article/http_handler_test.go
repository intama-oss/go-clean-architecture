package article

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go-clean-architecture/internal/domain"
	"go-clean-architecture/mocks"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHttpArticleHandler_Fetch(t *testing.T) {
	var mockArticle domain.Article
	var mockArticle2 domain.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	err = faker.FakeData(&mockArticle2)
	assert.NoError(t, err)
	mockService := new(mocks.ArticleService)
	mockListArticle := make([]*domain.Article, 0)
	mockListArticle = append(mockListArticle, &mockArticle, &mockArticle2)

	t.Run("success", func(t *testing.T) {
		size := uint(1)
		page := uint(1)
		mockService.On("Fetch", page, size, &domain.Article{}).
			Return(mockListArticle, uint(2), nil).Once()
		mockService.On("Count", &domain.Article{}).
			Return(int64(2), nil).Once()
		app := fiber.New()
		NewHttpHandler(app, mockService)
		resp, err := app.Test(httptest.NewRequest("GET", "/?page=1&size=1", nil))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "2", resp.Header.Get("X-Cursor"))
		assert.Equal(t, "2", resp.Header.Get("X-Total-Count"))
		assert.Equal(t, "2", resp.Header.Get("X-Max-Page"))
		mockService.AssertExpectations(t)
	})

	t.Run("success with search", func(t *testing.T) {
		size := uint(1)
		page := uint(1)
		mockService.On("Fetch", page, size, &domain.Article{Title: mockArticle.Title}).
			Return(mockListArticle, uint(2), nil).Once()
		mockService.On("Count", &domain.Article{Title: mockArticle.Title}).
			Return(int64(2), nil).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		resp, err := app.Test(httptest.NewRequest("GET", "/?page=1&size=1&q="+mockArticle.Title, nil))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(bodyBytes), mockArticle.Title)
		mockService.AssertExpectations(t)
	})

	t.Run("success with no data", func(t *testing.T) {
		mockNewService := new(mocks.ArticleService)
		mockNewService.On("Fetch", uint(1), uint(10), &domain.Article{}).
			Return(nil, uint(0), nil).Once()

		app := fiber.New()
		NewHttpHandler(app, mockNewService)
		resp, err := app.Test(httptest.NewRequest("GET", "/?page=1&size=10", nil))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(bodyBytes), "[]")
		mockNewService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockNewService := new(mocks.ArticleService)
		mockNewService.On("Fetch", uint(1), uint(10), &domain.Article{}).
			Return(nil, uint(0), errors.New("unexpected Error")).Once()

		app := fiber.New()
		NewHttpHandler(app, mockNewService)
		resp, err := app.Test(httptest.NewRequest("GET", "/?page=1&size=10", nil))
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		mockNewService.AssertExpectations(t)
	})

	t.Run("error total item", func(t *testing.T) {
		mockNewService := new(mocks.ArticleService)
		mockNewService.On("Fetch", uint(1), uint(10), &domain.Article{}).
			Return(mockListArticle, uint(2), nil).Once()
		mockNewService.On("Count", &domain.Article{}).
			Return(int64(0), errors.New("unexpected Error")).Once()

		app := fiber.New()
		NewHttpHandler(app, mockNewService)
		resp, err := app.Test(httptest.NewRequest("GET", "/?page=1&size=10", nil))
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		mockNewService.AssertExpectations(t)
	})
}

func TestHttpArticleHandler_Fetch_WithErrorSize(t *testing.T) {
	app := fiber.New()
	NewHttpHandler(app, nil)
	resp, err := app.Test(httptest.NewRequest("GET", "/?page=1&size=0", nil))
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestHttpArticleHandler_Fetch_WithErrorPage(t *testing.T) {
	app := fiber.New()
	NewHttpHandler(app, nil)
	resp, err := app.Test(httptest.NewRequest("GET", "/?page=0&size=1", nil))
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestHttpArticleHandler_GetByID(t *testing.T) {
	var mockArticle domain.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockService := new(mocks.ArticleService)

	t.Run("success", func(t *testing.T) {
		mockService.On("GetByID", mockArticle.ID).
			Return(&mockArticle, nil).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		id := strconv.Itoa(int(mockArticle.ID))
		resp, err := app.Test(httptest.NewRequest("GET", "/"+id, nil))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetByID", mockArticle.ID).
			Return(nil, fiber.ErrNotFound).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		id := strconv.Itoa(int(mockArticle.ID))
		resp, err := app.Test(httptest.NewRequest("GET", "/"+id, nil))
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error-parsing-id", func(t *testing.T) {
		app := fiber.New()
		NewHttpHandler(app, mockService)
		resp, err := app.Test(httptest.NewRequest("GET", "/abc", nil))
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("GetByID", mockArticle.ID).
			Return(nil, errors.New("unexpected Error")).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		id := strconv.Itoa(int(mockArticle.ID))
		resp, err := app.Test(httptest.NewRequest("GET", "/"+id, nil))
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestHttpArticleHandler_Store(t *testing.T) {
	var mockArticleStoreRequest domain.ArticleStoreRequest
	err := faker.FakeData(&mockArticleStoreRequest)
	assert.NoError(t, err)
	mockArticle := &domain.Article{
		Title:    mockArticleStoreRequest.Title,
		Content:  mockArticleStoreRequest.Content,
		AuthorID: mockArticleStoreRequest.AuthorID,
	}
	mockService := new(mocks.ArticleService)

	t.Run("success", func(t *testing.T) {
		mockService.On("Store", mockArticle).
			Return(nil).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		bodyRequest, err := json.Marshal(mockArticleStoreRequest)
		assert.NoError(t, err)
		bodyRequestIO := io.NopCloser(bytes.NewReader(bodyRequest))
		req := httptest.NewRequest("POST", "/", bodyRequestIO)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyRequest)))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("Store", mockArticle).
			Return(errors.New("unexpected Error")).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		bodyRequest, err := json.Marshal(mockArticleStoreRequest)
		assert.NoError(t, err)
		bodyRequestIO := io.NopCloser(bytes.NewReader(bodyRequest))
		req := httptest.NewRequest("POST", "/", bodyRequestIO)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyRequest)))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestHttpArticleHandler_Update(t *testing.T) {
	var mockArticleUpdateRequest domain.ArticleUpdateRequest
	err := faker.FakeData(&mockArticleUpdateRequest)
	assert.NoError(t, err)

	mockArticle := domain.Article{
		Title:   mockArticleUpdateRequest.Title,
		Content: mockArticleUpdateRequest.Content,
	}

	mockService := new(mocks.ArticleService)

	t.Run("success", func(t *testing.T) {
		mockService.On("Update", &mockArticle).
			Return(nil).
			Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		bodyRequest, err := json.Marshal(mockArticleUpdateRequest)
		assert.NoError(t, err)
		bodyRequestIO := io.NopCloser(bytes.NewReader(bodyRequest))
		req := httptest.NewRequest("PUT", "/"+strconv.Itoa(int(mockArticle.ID)), bodyRequestIO)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyRequest)))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error-parsing-id", func(t *testing.T) {
		app := fiber.New()
		NewHttpHandler(app, mockService)
		bodyRequest, err := json.Marshal(mockArticleUpdateRequest)
		assert.NoError(t, err)
		bodyRequestIO := io.NopCloser(bytes.NewReader(bodyRequest))
		req := httptest.NewRequest("PUT", "/abc", bodyRequestIO)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyRequest)))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("Update", &mockArticle).
			Return(errors.New("unexpected Error")).
			Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		bodyRequest, err := json.Marshal(mockArticleUpdateRequest)
		assert.NoError(t, err)
		bodyRequestIO := io.NopCloser(bytes.NewReader(bodyRequest))
		req := httptest.NewRequest("PUT", "/"+strconv.Itoa(int(mockArticle.ID)), bodyRequestIO)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyRequest)))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestHttpArticleHandler_Delete(t *testing.T) {
	var mockArticle domain.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockService := new(mocks.ArticleService)

	t.Run("success", func(t *testing.T) {
		mockService.On("Delete", mockArticle.ID).
			Return(nil).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		resp, err := app.Test(httptest.NewRequest("DELETE", "/"+strconv.Itoa(int(mockArticle.ID)), nil))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error-parsing-id", func(t *testing.T) {
		app := fiber.New()
		NewHttpHandler(app, mockService)
		resp, err := app.Test(httptest.NewRequest("DELETE", "/abc", nil))
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("Delete", mockArticle.ID).
			Return(errors.New("unexpected Error")).Once()

		app := fiber.New()
		NewHttpHandler(app, mockService)
		resp, err := app.Test(httptest.NewRequest("DELETE", "/"+strconv.Itoa(int(mockArticle.ID)), nil))
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
