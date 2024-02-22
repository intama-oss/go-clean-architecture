package xlogger

import (
	"github.com/rs/zerolog"
	"go-clean-architecture/internal/config"
	"os"
	"sync"
	"time"
)

var (
	Logger     *zerolog.Logger
	loggerOnce sync.Once
)

func Setup(cfg config.Config) {
	loggerOnce.Do(func() {
		if cfg.IsDevelopment {
			l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
			Logger = &l
			return
		}
		l := zerolog.New(os.Stderr).With().Timestamp().Logger()
		Logger = &l
	})
}
