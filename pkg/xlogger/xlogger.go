package xlogger

import (
	"github.com/rs/zerolog"
	"go-clean-architecture/internal/config"
	"os"
	"time"
)

var (
	Logger *zerolog.Logger
)

func Setup(cfg config.Config) {
	if cfg.IsDevelopment {
		l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
		Logger = &l
		return
	}
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	Logger = &l
}
