package handler

import (
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kewka/go-url-shortener/internal/handler/html"
	"github.com/kewka/go-url-shortener/internal/service"
)

//go:embed public
var publicFs embed.FS

type Config struct {
	Service   service.Service
	PublicUrl string
}

type handler struct {
	mux *chi.Mux
	cfg Config
}

func New(cfg Config) http.Handler {
	h := &handler{
		cfg: cfg,
		mux: chi.NewRouter(),
	}
	h.mux.Get("/", h.handleIndex())
	h.mux.Post("/", h.handleShortenUrl())
	h.mux.Get("/public/*", h.handlePublic())
	h.mux.Get("/url/{code}", h.handleUrl())
	h.mux.Get("/{code}", h.handleUrlRedirect())
	h.mux.NotFound(h.handleNotFound())
	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) handlePublic() http.HandlerFunc {
	return http.FileServer(http.FS(publicFs)).ServeHTTP
}

func (h *handler) handleUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := h.cfg.Service.FindUrlByCode(r.Context(), chi.URLParam(r, "code"))
		if err != nil {
			if err == service.ErrUrlNotFound {
				h.mux.NotFoundHandler().ServeHTTP(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		html.Url(w, html.UrlParams{
			Url:       u,
			PublicUrl: h.cfg.PublicUrl,
		})
	}
}

func (h *handler) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		html.NotFound(w)
	}
}

func (h *handler) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html.Index(w, html.IndexParams{})
	}
}

func (h *handler) handleUrlRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := h.cfg.Service.FindUrlByCode(r.Context(), chi.URLParam(r, "code"))
		if err != nil {
			if err == service.ErrUrlNotFound {
				h.mux.NotFoundHandler().ServeHTTP(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Redirect(w, r, u.Url, http.StatusMovedPermanently)
	}
}

func (h *handler) handleShortenUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rawurl := r.FormValue("url")
		u, err := h.cfg.Service.ShortenUrl(r.Context(), rawurl)
		if err != nil {
			if errors.Is(err, service.ErrUrlInvalid) {
				w.WriteHeader(http.StatusBadRequest)
				html.Index(w, html.IndexParams{
					ErrorMessage: err.Error(),
					Url:          rawurl,
				})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/url/%v", u.Code), http.StatusFound)
	}
}
