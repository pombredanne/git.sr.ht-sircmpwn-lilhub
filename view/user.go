package view

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"git.sr.ht/~sircmpwn/lilhub/github"
)

type UserPage struct {
	Page
	User *github.User
}

type OrgPage struct {
	Page
	Org *github.Organization
}

// GET /:login
func ProfileIndex(c echo.Context) error {
	ctx := c.Request().Context()
	client := github.ForContext(c)

	login := c.Param("login")
	user, org, _ := github.FetchProfile(client, ctx, login)
	if user != nil {
		title := user.Login
		if user.Name != nil {
			title = fmt.Sprintf("%s (%s)", user.Login, *user.Name)
		}
		return c.Render(http.StatusOK, "user.html", &UserPage{
			Page: NewPage(c, title),
			User: user,
		})
	} else if org != nil {
		title := org.Login
		if org.Name != nil {
			title = fmt.Sprintf("%s (%s)", org.Login, *org.Name)
		}
		return c.Render(http.StatusOK, "org.html", &OrgPage{
			Page: NewPage(c, title),
			Org:  org,
		})
	} else {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
}
