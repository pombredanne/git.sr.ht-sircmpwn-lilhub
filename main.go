package main

import (
	"html/template"
	"io"
	"net/http"
	"os"

	"git.sr.ht/~emersion/gqlclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/oauth2"

	"git.sr.ht/~sircmpwn/lilhub/github"
)

type Page struct {
	Title string
}

type UserPage struct {
	Page
	User *github.User
}

func main() {
	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} ${method} ${path} ${status} ${time_rfc3339} ${error}\n",
	}))
	e.Use(middleware.Recover())

	tok := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: os.Args[1],
		TokenType:   "bearer",
	})

	e.GET("/:user", func(c echo.Context) error {
		ctx := c.Request().Context()

		oclient := oauth2.NewClient(ctx, tok)
		client := gqlclient.New("https://api.github.com/graphql", oclient)

		username := c.Param("user")
		user, err := github.FetchUserIndex(client, ctx, username)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "user.html", &UserPage{
			Page: Page{
				Title: user.Login,
			},
			User: user,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
