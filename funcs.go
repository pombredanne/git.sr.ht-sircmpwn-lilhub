package main

import (
	"bytes"
	"html/template"
	"net/url"
	"path/filepath"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"git.sr.ht/~sircmpwn/lilhub/github"
)

func templateFuncs() template.FuncMap {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	return template.FuncMap{
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

		"rewriteurl": func(u github.URI) string {
			uri, err := url.Parse(string(u))
			if err != nil {
				panic(err)
			}
			uri.Scheme = ""
			uri.Host = ""
			return uri.String()
		},

		"date": func(ts *github.GitTimestamp) string {
			date, err := time.Parse("2006-01-02T15:04:05-07:00", string(*ts))
			if err != nil {
				panic(err)
			}
			return date.Format("Mon Jan 2 2006 at 15:04 PM")
		},

		"hl": func(text, filename string) template.HTML {
			lex := lexers.Match(filename)
			if lex == nil {
				lex = lexers.Fallback
			}
			style := styles.Get("github")
			formatter := html.New(
				// TODO: should not use standalone
				html.Standalone(true),
				html.WithLineNumbers(true),
				html.LinkableLineNumbers(true, "L"))
			iter, err := lex.Tokenise(nil, text)
			if err != nil {
				panic(err)
			}

			var buf bytes.Buffer
			err = formatter.Format(&buf, style, iter)
			if err != nil {
				panic(err)
			}
			return template.HTML(buf.String())
		},

		"parent": func(path []string) string {
			return filepath.Join(path[:len(path)-1]...)
		},
	}
}
