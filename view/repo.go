package view

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"

	"git.sr.ht/~sircmpwn/lilhub/github"
)

type RepoIndexPage struct {
	Page
	Repository *github.Repository
	Readme     *github.Repository
}

// GET /:owner/:repo
func RepoHome(c echo.Context) error {
	ctx := c.Request().Context()
	client := github.ForContext(c)

	owner := c.Param("owner")
	reponame := c.Param("repo")
	repo, _ := github.FetchRepoIndex(client, ctx, owner, reponame)
	// XXX: Errors are ignored, need more general solution

	// TODO: this probably fails if there is no latest tree or something
	commit := repo.Object.Value.(*github.Commit)
	entries := commit.Tree.Entries

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type > entries[j].Type
		}
		return entries[i].Name < entries[j].Name
	})

	readme, _ := github.FetchREADME(client, ctx, owner, reponame)

	return c.Render(http.StatusOK, "repo-index.html", &RepoIndexPage{
		Page: NewPage(c, fmt.Sprintf("%s/%s",
			repo.Owner.Login,
			repo.Name)),
		Repository: repo,
		Readme:     readme,
	})
}
