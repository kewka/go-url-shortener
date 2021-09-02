package test

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kewka/go-url-shortener/internal/handler"
	"github.com/kewka/go-url-shortener/internal/postgres"
	"github.com/kewka/go-url-shortener/internal/service"
)

var PublicUrl = "https://kewka.sh/"

type Env struct {
	Handler http.Handler
	Dbpool  *pgxpool.Pool
	Service service.Service
}

func NewEnv() (*Env, error) {
	ret := &Env{}
	var err error
	ret.Dbpool, err = postgres.Setup(context.Background())
	if err != nil {
		return ret, err
	}
	ret.Service = service.New(ret.Dbpool)
	ret.Handler = handler.New(handler.Config{
		Service:   ret.Service,
		PublicUrl: PublicUrl,
	})
	return ret, nil
}

func (env *Env) Cleanup() {
	env.Dbpool.Close()
}
