package view

import (
	"github.com/labstack/echo/v4"
)

type Page struct {
	Title     string
	GitHubURL string
}

func NewPage(c echo.Context, title string) Page {
	ghurl := *c.Request().URL
	ghurl.Scheme = "https"
	ghurl.Host = "github.com"

	return Page{
		Title:     title,
		GitHubURL: ghurl.String(),
	}
}
