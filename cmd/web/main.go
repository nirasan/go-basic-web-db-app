package main

import (
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nirasan/go-basic-web-db-app/app"
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

	e.GET("/", app.IndexHandler)
	e.GET("/create", app.CreateHandler)
	e.POST("/create", app.CreateHandler)
	e.GET("/update/:id", app.UpdateHandler)
	e.POST("/update/:id", app.UpdateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		e.Logger.Printf("Defaulting to port %s", port)
	}
	e.Logger.Fatal(e.Start(":" + port))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
