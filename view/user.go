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

func UserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	client := github.ForContext(c)

	username := c.Param("user")
	user, _ := github.FetchUserIndex(client, ctx, username)
	// XXX: Errors are ignored, need more general solution

	return c.Render(http.StatusOK, "user.html", &UserPage{
		Page: NewPage(c, fmt.Sprintf("%s (%s)", user.Login, *user.Name)),
		User: user,
	})
}
