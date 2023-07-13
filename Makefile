STYLES=\
	css/blob.scss \
	css/layout.scss \
	css/main.scss \
	css/markdown.scss \
	css/nav.scss \
	css/repo.scss \
	css/typography.scss \
	css/user-profile.scss \
	css/variables.scss \
	css/well.scss

static/css/main.css: $(STYLES)
	mkdir -p static/css
	sassc --style compressed css/main.scss $@

run: static/css/main.css
	go run .

watch:
	while inotifywait -r -e close_write css || true; \
	do \
		make static/css/main.css; \
	done

.PHONY: run
