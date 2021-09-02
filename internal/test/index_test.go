package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/kewka/go-url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestIndexPage(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	env.Handler.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	assert.Nil(t, err)
	formEl := doc.Find("form")
	assert.Equal(t, "POST", formEl.AttrOr("method", ""))
	assert.Equal(t, "/", formEl.AttrOr("action", ""))
}

func TestIndexPageInvalid(t *testing.T) {
	t.Parallel()
	form := url.Values{}
	form.Set("url", "invalid_url")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	env.Handler.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, form.Get("url"), doc.Find(`input[name="url"]`).AttrOr("value", ""))
	assert.Contains(t, doc.Find(`.Typography--error`).Text(), fmt.Sprintf("Error: %v", service.ErrUrlInvalid))
}

func TestIndexPageSuccess(t *testing.T) {
	t.Parallel()
	form := url.Values{}
	form.Set("url", "https://github.com/kewka")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	env.Handler.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusFound, res.StatusCode)
	assert.True(t, strings.HasPrefix(res.Header.Get("Location"), "/url/"))
}
