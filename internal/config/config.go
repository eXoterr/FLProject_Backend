package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	API      API      `yaml:"api"`
	Env      Env      `yaml:"env"`
	Logger   Logger   `yaml:"logger"`
	Database Database `yaml:"database"`
	Security Security `yaml:"security"`
}

type API struct {
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	ListenAddr  string        `yaml:"listen_addr" env-default:"127.0.0.1:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"4s"`
	CORS        CORS          `yaml:"cors"`
}

type CORS struct {
	Origins     []string `yaml:"origins" end-required:"true"`
	Methods     []string `yaml:"methods" end-required:"true"`
	Headers     []string `yaml:"headers" end-required:"true"`
	Credentials bool     `yaml:"credentials" env-required:"true"`
}

type Env struct {
	Mode string `yaml:"mode" env-required:"true"`
}

type Logger struct {
	LogLevel string `yaml:"log_level" env-default:"info"`
	Format   string `yaml:"format" env-default:"text"`
}

type Security struct {
	BCryptCost int `yaml:"bcrypt_cost" env-required:"true"`
	JWT        JWT
}

type JWT struct {
	RefreshTokenTTL int64  `yaml:"refresh_ttl" env-required:"true"`
	AccessTokenTTL  int64  `yaml:"access_ttl" env-required:"true"`
	Key             string `env:"JWT_KEY" env-required:"true"`
}

type Database struct {
	SQLite     SQLite     `yaml:"sqlite"`
	PostgreSQL PostgreSQL `yaml:"postgre"`
}

type SQLite struct {
	Path string `yaml:"path" env-default:"store.sqlite3"`
}

type PostgreSQL struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `env:"DB_PASS" env-required:"true"`
	Database string `yaml:"database" env-default:"postgres"`
}

func MustLoad(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config path %s does not exist: %s", path, err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatalf("unable to read config file: %s", err)
	}

	err = cfg.validate()
	if err != nil {
		log.Fatalf("unable to validate config values: %s", err)
	}

	return &cfg
}
