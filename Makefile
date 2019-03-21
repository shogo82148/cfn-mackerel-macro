SRC_FILES=$(shell find . -type f -name '*.go')

help: ## Show this text.
	# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all test clean help

all: resource/main ## Build a package

resource/main: resource $(SRC_FILES) go.mod go.sum
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o resource/main .

resource:
	mkdir -p resource

test:
	go test -v -race ./...

clean:
	@rm -f packaged.yaml
