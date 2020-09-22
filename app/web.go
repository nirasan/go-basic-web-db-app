package app

import (
	"fmt"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
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

func DeleteHandler(c echo.Context) error {
	if c.Request().Method != "POST" {
		return fmt.Errorf("invalid method")
	}

	client, err := NewDBClient()
	if err != nil {
		return err
	}

	res := client.books.Find(c.Param("id"))
	var book Book
	if err := res.One(&book); err != nil {
		return err
	}

	if err := client.books.Find(book.ID).Delete(); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func LoginHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", struct {
	}{})
}

func LoginSuccessHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "login_success.html", struct {
		}{})
	}

	ctx := c.Request().Context()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}

	// get token
	idToken := c.FormValue("token")
	expiresIn := time.Hour * 24 * 5

	// verify token
	decoded, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return err
	}
	// Return error if the sign-in is older than 5 minutes.
	if time.Now().Unix()-int64(decoded.Claims["auth_time"].(float64)) > 5*60 {
		return fmt.Errorf("Recent sign-in required")
	}

	// set token to cookie
	sessionCookie, err := client.SessionCookie(ctx, idToken, expiresIn)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = sessionCookie
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	// create User if not exists
	dbclient, err := NewDBClient()
	if err != nil {
		return err
	}
	ok, err := dbclient.users.Find("id", decoded.UID).Exists()
	if err != nil {
		return err
	}
	if !ok {
		_, err := dbclient.users.Insert(&User{
			ID:   decoded.UID,
			Name: c.FormValue("username"),
		})
		if err != nil {
			return err
		}
	}

	return c.Redirect(http.StatusFound, "/user")
}

func LogoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.MaxAge = 0
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/user")
}

func UserHandler(c echo.Context) error {
	ctx := c.Request().Context()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}

	cookie, err := c.Cookie("session")
	if err != nil {
		return err
	}
	decoded, err := client.VerifySessionCookieAndCheckRevoked(ctx, cookie.Value)
	if err != nil {
		return err
	}

	dbclient, err := NewDBClient()
	if err != nil {
		return err
	}
	var user User
	if err := dbclient.users.Find("id", decoded.UID).One(&user); err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("%+v", user))
}
