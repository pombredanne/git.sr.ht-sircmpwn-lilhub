static/css/main.css: css/main.scss
	mkdir -p static/css
	sassc --style compressed $< $@

run: static/css/main.css
	go run main.go

watch:
	while inotifywait -e close_write css/main.scss || true; \
	do \
		make static/css/main.css; \
	done

.PHONY: run
