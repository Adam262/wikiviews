.PHONY: all
all: test

.PHONY: run
run:
	air go run main.go

.PHONY: test
test: vet
	go test ./...

.PHONY: vet
vet:
	go fmt ./...
	go vet ./...

.PHONY: build 
build:
	docker build . --tag wikiviews

.PHONY: docker-run 
docker-run: build
	docker run --name wikiviews -p 8080:8080 -d wikiviews 

.PHONY: mod
mod:
	go mod tidy
	go mod verify
	go mod vendor
