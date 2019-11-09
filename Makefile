SRC_FILES=$(shell find . -type f -name '*.go')

help: ## Show this text.
	# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all test clean help release

all: macro.zip resource.zip ## Build a package

resource/main: $(SRC_FILES) go.mod go.sum
	mkdir -p resource
	./run-in-docker.sh go build -o resource/main .

macro.zip: macro/app.py
	cd macro && zip -r ../macro.zip .

resource.zip: resource/main
	cd resource && zip -r ../resource.zip .

test:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	cfn-lint --override-spec cfn-resource-specification.json example.yaml

clean:
	@rm -f packaged.yaml
