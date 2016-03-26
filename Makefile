GO_FILES := $(shell find . -type f -name '*.go' -not -path "./Godeps/*" -not -path "./vendor/*")
GO_PACKAGES := $(shell go list ./... | sed "s/github.com\/heroku\/cytokine/./" | grep -v "^./vendor/")

build:
	go build -v $(GO_PACKAGES)

test: build
	go fmt $(GO_PACKAGES)
	go test -race -i $(GO_PACKAGES)
	go test -race -v $(GO_PACKAGES)
