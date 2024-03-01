.PHONY: run build docker-run
run:
	air go run main.go

build:
	docker build . --tag wikiviews

docker-run: build
	docker run --name wikiviews -p 8080:8080 -d wikiviews 