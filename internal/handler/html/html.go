package html

import (
	"embed"
	"html/template"
	"io"

	"github.com/kewka/go-url-shortener/internal/proto"
)

//go:embed *.html
var fs embed.FS

var (
	indexTmpl    = parse("index.html")
	urlTmpl      = parse("url.html")
	notFoundTmpl = parse("not-found.html")
)

func parse(file string) *template.Template {
	return template.Must(template.ParseFS(fs, "layout.html", file))
}

type IndexParams struct {
	ErrorMessage string
	Url          string
}

func Index(w io.Writer, params IndexParams) {
	indexTmpl.Execute(w, params)
}

type UrlParams struct {
	Url       proto.Url
	PublicUrl string
}

func Url(w io.Writer, params UrlParams) {
	urlTmpl.Execute(w, params)
}

func NotFound(w io.Writer) {
	notFoundTmpl.Execute(w, nil)
}
