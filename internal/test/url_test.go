package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestUrlRedirect(t *testing.T) {
	t.Parallel()
	rawurl := "https://github.com/kewka"
	u, err := env.Service.ShortenUrl(context.Background(), rawurl)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%v", u.Code), nil)
	env.Handler.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)
	assert.Equal(t, rawurl, res.Header.Get("Location"))
}

func TestUrlPage(t *testing.T) {
	t.Parallel()
	rawurl := "https://github.com/kewka"
	u, err := env.Service.ShortenUrl(context.Background(), rawurl)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/url/%v", u.Code), nil)
	env.Handler.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, PublicUrl+u.Code, doc.Find(`.UrlPage-url-input`).AttrOr("value", ""))
	assert.Equal(t, u.Url, doc.Find(`.UrlPage-long-url a`).AttrOr("href", ""))
}
