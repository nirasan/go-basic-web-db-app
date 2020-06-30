package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	e.Static("/static", "static")
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t

	e.GET("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		e.Logger.Printf("Defaulting to port %s", port)
	}
	e.Logger.Fatal(e.Start(":"+port))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func indexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", struct{
		Title string
	}{
		Title: "my home page",
	})
}