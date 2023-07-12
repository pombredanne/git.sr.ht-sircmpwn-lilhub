//go:build generate
// +build generate

package github

import (
	_ "git.sr.ht/~emersion/gqlclient/cmd/gqlclientgen"
)

//go:generate go run git.sr.ht/~emersion/gqlclient/cmd/gqlclientgen -s schema.graphqls -q queries.graphql -o gql.go -n github
