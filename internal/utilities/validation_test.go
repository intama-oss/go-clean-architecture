package utilities

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestExtractStructFromValidator(t *testing.T) {
	type Payload struct {
		Name string `json:"name" validate:"required"`
	}

	payload := Payload{Name: "John Doe"}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("parser", &payload)
		return c.Next()
	})

	app.Post("/", func(c *fiber.Ctx) error {
		p := ExtractStructFromValidator[Payload](c)
		return c.JSON(p)
	})

	req, err := app.Test(httptest.NewRequest("POST", "/", nil))
	assert.NoError(t, err)
	assert.Equal(t, 200, req.StatusCode)
}

func TestExtractStructFromValidator_NotFound(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		p := ExtractStructFromValidator[struct{}](c)
		return c.JSON(p)
	})

	req, err := app.Test(httptest.NewRequest("POST", "/", nil))
	assert.NoError(t, err)
	assert.Equal(t, 200, req.StatusCode)
}
