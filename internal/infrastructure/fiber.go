package infrastructure

import (
	"fmt"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go-clean-architecture/internal/article"
	"go-clean-architecture/internal/docs"
	"go-clean-architecture/pkg/xlogger"
)

func Run() {
	logger := xlogger.Logger

	app := fiber.New(fiber.Config{
		ProxyHeader:           cfg.ProxyHeader,
		DisableStartupMessage: true,
		ErrorHandler:          defaultErrorHandler,
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
		Fields: cfg.LogFields,
	}))
	app.Use(recover2.New())
	app.Use(etag.New())
	app.Use(requestid.New())

	api := app.Group("/api")
	docs.NewHttpHandler(api.Group("/docs"))
	article.NewHttpHandler(api.Group("/articles"), articleService)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info().Msgf("Server is running on address: %s", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
