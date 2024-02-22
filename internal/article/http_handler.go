package article

import (
	"github.com/gofiber/fiber/v2"
	"go-clean-architecture/internal/domain"
	"go-clean-architecture/internal/middleware/validation"
	"go-clean-architecture/internal/utilities"
	"strconv"
)

type HttpArticleHandler struct {
	articleSvc domain.ArticleService
}

func NewHttpHandler(r fiber.Router, articleSvc domain.ArticleService) {
	handler := &HttpArticleHandler{
		articleSvc: articleSvc,
	}
	r.Post("/", validation.New[domain.ArticleStoreRequest](), handler.Store)
	r.Get("/", handler.Fetch)
	r.Get("/:id", handler.GetByID)
	r.Put("/:id", validation.New[domain.ArticleUpdateRequest](), handler.Update)
	r.Delete("/:id", handler.Delete)
}

// Fetch used to get list of articles
//
//	@Summary		Get list of articles
//	@Description	Get list of articles
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"Page number (default 1)"
//	@Param			size	query		int				false	"Size of page (default 10)"
//	@Param			q		query		string			false	"Search query"
//	@Header			200		{string}	X-Cursor		"Next page"
//	@Header			200		{string}	X-Total-Count	"Total item"
//	@Header			200		{string}	X-Max-Page		"Max page"
//	@Success		200		{array}		domain.Article	"List of articles"
//	@Failure		400		{object}	domain.Error	"Bad Request"
//	@Failure		500		{object}	domain.Error	"Internal Server Error"
//	@Router			/articles [get]
func (h *HttpArticleHandler) Fetch(c *fiber.Ctx) error {
	page, size, query := c.QueryInt("page", 1), c.QueryInt("size", 10), c.Query("q")
	if page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "page must be a positive integer",
		})
	}
	if size <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "size must be a positive integer",
		})
	}

	filter := &domain.Article{Title: query}
	articles, nextPage, err := h.articleSvc.Fetch(uint(page), uint(size), filter)
	if err != nil {
		return err
	}

	if articles == nil {
		return c.JSON([]domain.Article{})
	}

	totalItem, err := h.articleSvc.Count(filter)
	if err != nil {
		return err
	}

	maxPage := int(totalItem) / size

	if nextPage > 0 && nextPage <= uint(maxPage) {
		c.Set("X-Cursor", strconv.Itoa(int(nextPage)))
	}
	c.Set("X-Total-Count", strconv.Itoa(int(totalItem)))
	c.Set("X-Max-Page", strconv.Itoa(maxPage))
	return c.JSON(articles)
}

// GetByID used to get article by id
//
//	@Summary		Get article by id
//	@Description	Get article by id
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Article ID"
//	@Success		200	{object}	domain.Article	"Article detail"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		404	{object}	domain.Error	"Not Found"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/articles/{id} [get]
func (h *HttpArticleHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	article, err := h.articleSvc.GetByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(article)
}

// Store used to store article
//
//	@Summary		Store article
//	@Description	Store article
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			article	body		domain.ArticleStoreRequest	true	"Article data"
//	@Success		201		{object}	domain.Article				"Article detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/articles [post]
func (h *HttpArticleHandler) Store(c *fiber.Ctx) error {
	articleReq := utilities.ExtractStructFromValidator[domain.ArticleStoreRequest](c)

	article := &domain.Article{
		Title:    articleReq.Title,
		Content:  articleReq.Content,
		AuthorID: articleReq.AuthorID,
	}

	if err := h.articleSvc.Store(article); err != nil {
		return err
	}

	c.Status(fiber.StatusCreated)
	return c.JSON(article)
}

// Update used to update article
//
//	@Summary		Update article
//	@Description	Update article
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Article ID"
//	@Param			article	body		domain.ArticleUpdateRequest	true	"Article data"
//	@Success		200		{object}	domain.Article				"Article detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		404		{object}	domain.Error				"Not Found"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/articles/{id} [put]
func (h *HttpArticleHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	articleReq := utilities.ExtractStructFromValidator[domain.ArticleUpdateRequest](c)

	article := &domain.Article{
		ID:      uint(id),
		Title:   articleReq.Title,
		Content: articleReq.Content,
	}

	if err := h.articleSvc.Update(article); err != nil {
		return err
	}

	return c.JSON(article)
}

// Delete used to delete article
//
//	@Summary		Delete article
//	@Description	Delete article
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Article ID"
//	@Success		200	{object}	domain.Message	"Success delete article"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/articles/{id} [delete]
func (h *HttpArticleHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.articleSvc.Delete(uint(id)); err != nil {
		return err
	}

	return c.JSON(domain.Message{
		Code:    fiber.StatusOK,
		Message: "Success delete article",
	})
}
