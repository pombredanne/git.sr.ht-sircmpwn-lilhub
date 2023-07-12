package main

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
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
	}
}
