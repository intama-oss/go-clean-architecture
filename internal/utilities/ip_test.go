package utilities

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func startTestServerWithPort(t *testing.T, beforeStarting func(app *fiber.App)) string {
	t.Helper()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	if beforeStarting != nil {
		beforeStarting(app)
	}

	ln, err := net.Listen(fiber.NetworkTCP4, "127.0.0.1:0")
	assert.NoError(t, err)

	go func() {
		assert.NoError(t, app.Listener(ln))
	}()

	time.Sleep(2 * time.Second)
	addr := ln.Addr().String()

	return addr
}

func TestGetPublicIP(t *testing.T) {
	t.Parallel()

	addr := startTestServerWithPort(t, func(app *fiber.App) {
		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(map[string]string{
				"ip": "103.186.202.7",
			})
		})
	})

	ip, err := GetPublicIP("http://" + addr)
	assert.NoError(t, err)
	assert.Equal(t, "103.186.202.7", ip)
}

func TestGetPublicIP_Error(t *testing.T) {
	t.Parallel()

	addr := startTestServerWithPort(t, func(app *fiber.App) {
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendStatus(500)
		})
	})

	_, err := GetPublicIP("http://" + addr)
	assert.Error(t, err)
}

func TestGetPublicIP_InvalidResponse(t *testing.T) {
	t.Parallel()

	addr := startTestServerWithPort(t, func(app *fiber.App) {
		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(map[string]interface{}{
				"ip": 123,
			})
		})
	})

	_, err := GetPublicIP("http://" + addr)
	assert.Error(t, err)
}

func TestGetPublicIP_EmptyResponse(t *testing.T) {
	t.Parallel()

	addr := startTestServerWithPort(t, func(app *fiber.App) {
		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(map[string]interface{}{
				"ip": "",
			})
		})
	})

	_, err := GetPublicIP("http://" + addr)
	assert.Error(t, err)
}

func TestGetPublicIP_ErrAgent(t *testing.T) {
	t.Parallel()

	_, err := GetPublicIP("")
	assert.Error(t, err)
}

func TestGetPublicIP_ErrAgentParse(t *testing.T) {
	t.Parallel()

	_, err := GetPublicIP("invalid://test")
	assert.Error(t, err)
}
