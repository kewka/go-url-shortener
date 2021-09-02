package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Database string `envconfig:"POSTGRES_DB" required:"true"`
	SSLMode  string `envconfig:"POSTGRES_SSLMODE" default:"disable"`
}

func (c Config) ConnectionString() string {
	return fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=%v",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

func LoadConfig() (Config, error) {
	ret := Config{}
	return ret, envconfig.Process("", &ret)
}

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, cfg.ConnectionString())
}

func Setup(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return NewPool(ctx, cfg)
}
