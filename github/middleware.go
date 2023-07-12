package github

import (
	"errors"

	"git.sr.ht/~emersion/gqlclient"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

const GITHUB_API_URL = "https://api.github.com/graphql"

// Creates middleware for GitHub clients with the provided default OAuth token.
func Middleware(token string) echo.MiddlewareFunc {
	tok := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		TokenType:   "bearer",
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			oclient := oauth2.NewClient(ctx, tok)
			client := gqlclient.New(GITHUB_API_URL, oclient)
			c.Set("github", client)
			return next(c)
		}
	}
}

// Returns a GitHub GraphQL client for the given context.
func ForContext(ctx echo.Context) *gqlclient.Client {
	client, ok := ctx.Get("github").(*gqlclient.Client)
	if !ok {
		panic(errors.New("Invalid GitHub context"))
	}
	return client
}
