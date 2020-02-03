package provider

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Provider struct {
	port int
	server *echo.Echo
}

type Template struct {
	templates *template.Template
}
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewProvider(port int) *Provider {
	return &Provider{port: port, server: echo.New()}
}

func (provider *Provider) RegisterSamplePage() {
	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}

	provider.server.Renderer = t
	provider.RegisterGetMethod("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	  })
}

func (provider *Provider) UseCORS() {
	provider.server.Use(middleware.CORS())
}

func (provider *Provider) RegisterPostMethod(path string, handler echo.HandlerFunc) {
	provider.server.POST(path, handler)
}

func (provider *Provider) RegisterGetMethod(path string, handler echo.HandlerFunc) {
	provider.server.GET(path, handler)
}

func (provider *Provider) Run() error {
	return provider.server.Start(fmt.Sprintf(":%d", provider.port))
}