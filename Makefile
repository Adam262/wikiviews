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

.PHONY: mod
mod:
	go mod tidy
	go mod verify
	go mod vendor

.PHONY: build
build: vet test mod
	go build -o ./cmd/http-server ./...

.PHONY: docker-build
docker-build: vet test mod
	docker build . --tag wikiviews --progress=plain
