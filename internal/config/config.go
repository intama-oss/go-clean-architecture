package config

type Config struct {
	Host          string   `env:"HOST"`
	Port          int      `env:"PORT" envDefault:"3000"`
	IsDevelopment bool     `env:"IS_DEVELOPMENT"`
	ProxyHeader   string   `env:"PROXY_HEADER"`
	Database      Database `envPrefix:"DB_"`
	LogFields     []string `env:"LOG_FIELDS" envSeparator:","`
}

type Database struct {
	Driver string `env:"DRIVER" envDefault:"sqlite"`
	DSN    string `env:"DSN" envDefault:"file::memory:?cache=shared"`
}
