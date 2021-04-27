package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jxskiss/base62"
	"github.com/kewka/go-url-shortener/model"
)

var (
	ErrUrlInvalid  = errors.New("invalid url")
	ErrUrlNotFound = errors.New("url not found")
)

type Service interface {
	ShortenUrl(ctx context.Context, rawurl string) (model.Url, error)
	FindUrlByCode(ctx context.Context, code string) (model.Url, error)
}

type service struct {
	dbpool *pgxpool.Pool
}

func New(dbpool *pgxpool.Pool) Service {
	return &service{
		dbpool: dbpool,
	}
}

func (svc *service) ShortenUrl(ctx context.Context, rawurl string) (model.Url, error) {
	ret := model.Url{}
	if _, err := url.ParseRequestURI(rawurl); err != nil {
		return ret, ErrUrlInvalid
	}
	tx, err := svc.dbpool.Begin(ctx)
	if err != nil {
		return ret, err
	}
	defer tx.Rollback(ctx)
	var id int64
	if err := tx.
		QueryRow(
			ctx,
			"INSERT INTO urls (code, url) VALUES ('', $1) RETURNING id",
			rawurl,
		).
		Scan(&id); err != nil {
		return ret, err
	}
	if err := tx.QueryRow(
		ctx,
		"UPDATE urls SET code = $1 WHERE id = $2 RETURNING id, code, url",
		svc.encode(id),
		id,
	).Scan(&ret.Id, &ret.Code, &ret.Url); err != nil {
		return ret, err
	}
	return ret, tx.Commit(ctx)
}

func (svc *service) encode(id int64) string {
	return string(base62.FormatInt(id))
}

func (svc *service) FindUrlByCode(ctx context.Context, code string) (model.Url, error) {
	ret := model.Url{}
	if err := svc.dbpool.QueryRow(
		ctx,
		"SELECT id, code, url FROM urls WHERE code = $1 LIMIT 1",
		code,
	).Scan(&ret.Id, &ret.Code, &ret.Url); err != nil {
		if err == pgx.ErrNoRows {
			return ret, ErrUrlNotFound
		}
		return ret, err
	}
	return ret, nil
}
