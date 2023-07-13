package view

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"

	"git.sr.ht/~sircmpwn/lilhub/github"
)

type RepoIndexPage struct {
	Page
	Repository *github.Repository
	Readme     *github.Repository
	Tree       TreeDetails
}

type TreeDetails struct {
	Path       []string
	Repository *github.Repository
	Tree       *github.Tree
	Ref        string
}

func sortTree(tree *github.Tree) {
	entries := tree.Entries
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type > entries[j].Type
		}
		return entries[i].Name < entries[j].Name
	})
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
	sortTree(commit.Tree)

	readme, _ := github.FetchREADME(client, ctx, owner, reponame)

	return c.Render(http.StatusOK, "repo-index.html", &RepoIndexPage{
		Page: NewPage(c, fmt.Sprintf("%s/%s",
			repo.Owner.Login,
			repo.Name)),
		Repository: repo,
		Readme:     readme,

		Tree: TreeDetails{
			Path:       nil,
			Repository: repo,
			Tree:       commit.Tree,
			Ref:        repo.DefaultBranchRef.Name,
		},
	})
}

type RepoBlobPage struct {
	Page
	Repository *github.Repository
	Blob       *github.Blob
	Ref        string
	Path       string
	Name       string
	Size       string
	Breadcrumb []string
}

// GET /:owner/:repo/blob/:ref/:path
func RepoBlob(c echo.Context) error {
	ctx := c.Request().Context()
	client := github.ForContext(c)

	owner := c.Param("owner")
	reponame := c.Param("repo")
	ref := c.Param("ref")
	filepath := c.Param("path")

	revspec := fmt.Sprintf("%s:%s", ref, filepath)
	repo, _ := github.FetchBlob(client, ctx, owner, reponame, revspec)

	dir, name := path.Split(filepath)
	breadcrumb := strings.Split(dir, "/")
	if strings.HasSuffix(dir, "/") {
		breadcrumb = breadcrumb[:len(breadcrumb)-1]
	}

	size := humanize.Bytes(uint64(repo.Object.Value.(*github.Blob).ByteSize))

	return c.Render(http.StatusOK, "repo-blob.html", &RepoBlobPage{
		Page: NewPage(c, fmt.Sprintf("%s/%s",
			repo.Owner.Login,
			repo.Name)),
		Repository: repo,
		Blob:       repo.Object.Value.(*github.Blob),
		Ref:        ref,
		Path:       filepath,
		Name:       name,
		Size:       size,
		Breadcrumb: breadcrumb,
	})
}

type RepoTreePage struct {
	Page
	Repository *github.Repository
	Tree       TreeDetails
}

// GET /:owner/:repo/tree/:ref/:path
func RepoTree(c echo.Context) error {
	ctx := c.Request().Context()
	client := github.ForContext(c)

	owner := c.Param("owner")
	reponame := c.Param("repo")
	ref := c.Param("ref")
	path := c.Param("path")

	revspec := fmt.Sprintf("%s:%s", ref, path)
	repo, _ := github.FetchTree(client, ctx, owner, reponame, revspec)
	tree := repo.Object.Value.(*github.Tree)
	sortTree(tree)

	return c.Render(http.StatusOK, "repo-tree.html", &RepoTreePage{
		Page: NewPage(c, fmt.Sprintf("%s/%s",
			repo.Owner.Login,
			repo.Name)),
		Repository: repo,

		Tree: TreeDetails{
			Path:       strings.Split(path, "/"),
			Repository: repo,
			Tree:       tree,
			Ref:        ref,
		},
	})
}
