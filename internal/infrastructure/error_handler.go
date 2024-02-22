package infrastructure

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-clean-architecture/internal/domain"
)

func defaultErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := fiber.ErrInternalServerError.Error()

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		msg = e.Message
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(code).JSON(&domain.Error{
		Code:    code,
		Message: msg,
	})
}
