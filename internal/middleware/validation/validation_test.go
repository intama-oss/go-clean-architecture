package validation

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestValidation_SuccessPost(t *testing.T) {
	type Payload struct {
		Name string `json:"name" validate:"required"`
	}

	app := fiber.New()
	app.Post("/", New[Payload](), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	body := `{"name":"John Doe"}`
	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestValidation_EmptyBodyPost(t *testing.T) {
	type Payload struct {
		Name string `json:"name" validate:"required"`
	}

	app := fiber.New()
	app.Post("/", New[Payload](), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	resp, err := app.Test(httptest.NewRequest("POST", "/", nil))
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestValidation_InvalidBodyPost(t *testing.T) {
	type Payload struct {
		Name string `json:"name" validate:"required,min=5"`
	}

	app := fiber.New()
	app.Use(New[Payload]())
	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	body := `{"name":"John"}`
	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
