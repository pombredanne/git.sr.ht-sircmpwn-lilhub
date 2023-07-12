package main

import (
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"git.sr.ht/~sircmpwn/lilhub/github"
	"git.sr.ht/~sircmpwn/lilhub/view"
)

func main() {
	e := echo.New()
	e.Static("/static", "static")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} ${method} ${path} ${status} ${time_rfc3339} ${error}\n",
	}))
	//e.Use(middleware.Recover())

	e.Renderer = &Template{
		templates: template.Must(template.
			New("templates").
			Funcs(templateFuncs()).
			ParseGlob("templates/*.html")),
	}

	e.Use(github.Middleware(os.Getenv("GITHUB_TOKEN")))

	e.GET("/:user", view.UserProfile)
	e.GET("/:owner/:repo", view.RepoHome)
	e.GET("/:owner/:repo/blob/:ref/:path", view.RepoBlob)
	e.GET("/:owner/:repo/tree/:ref/:path", view.RepoTree)

	e.Logger.Fatal(e.Start(":1323"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
