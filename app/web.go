package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	client, err := NewDBClient()
	if err != nil {
		return err
	}

	var books []Book
	if err := client.books.Find().All(&books); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index.html", struct {
		Title string
		Books []Book
	}{
		Title: "list of books",
		Books: books,
	})
}

func CreateHandler(c echo.Context) error {
	if c.Request().Method == "GET" {
		return c.Render(http.StatusOK, "form.html", struct {
			Title string
			Book  *Book
		}{
			Title: "create book",
			Book:  nil,
		})
	}

	b := new(Book)
	if err := c.Bind(b); err != nil {
		return err
	}

	client, err := NewDBClient()
	if err != nil {
		return err
	}

	_, err = client.books.Insert(b)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func UpdateHandler(c echo.Context) error {
	client, err := NewDBClient()
	if err != nil {
		return err
	}

	if c.Request().Method == "GET" {
		res := client.books.Find(c.Param("id"))
		var book Book
		if err := res.One(&book); err != nil {
			return err
		}

		return c.Render(http.StatusOK, "form.html", struct {
			Title string
			Book  *Book
		}{
			Title: "create book",
			Book:  &book,
		})
	}

	book := new(Book)
	if err := c.Bind(book); err != nil {
		return err
	}

	if err := client.books.Find(book.ID).Update(book); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
