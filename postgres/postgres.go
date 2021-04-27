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
		"postgresql://%v:%v@%v:%v/%v?sslmode=%v",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

func LoadConfig() (Config, error) {
	ret := Config{}
	if err := envconfig.Process("", &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, cfg.ConnectionString())
}
