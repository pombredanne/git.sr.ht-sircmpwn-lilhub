package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

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
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	funcs := template.FuncMap{
		"md": func(text string) template.HTML {
			var buf bytes.Buffer
			buf.WriteString(`<div class="markdown">`)
			err := md.Convert([]byte(text), &buf)
			if err != nil {
				return template.HTML(
					template.HTMLEscapeString(err.Error()))
			}
			buf.WriteString(`</div>`)
			return template.HTML(buf.String())
		},
		"html": func(text string) template.HTML {
			return template.HTML(text)
		},
	}

	e := echo.New()
	e.Static("/static", "static")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} ${method} ${path} ${status} ${time_rfc3339} ${error}\n",
	}))
	e.Use(middleware.Recover())

	e.Renderer = &Template{
		templates: template.Must(template.
			New("templates").
			Funcs(funcs).
			ParseGlob("templates/*.html")),
	}

	e.Use(github.Middleware(os.Getenv("GITHUB_TOKEN")))

	e.GET("/:user", func(c echo.Context) error {
		ctx := c.Request().Context()
		client := github.ForContext(c)

		username := c.Param("user")
		user, _ := github.FetchUserIndex(client, ctx, username)
		// XXX: Errors are ignored, need more general solution

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
