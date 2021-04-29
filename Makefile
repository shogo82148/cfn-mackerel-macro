SRC_FILES=$(shell find . -type f -name '*.go')

help: ## Show this text.
	# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all test clean help release

all: macro.zip resource.zip template.yaml ## Build a package

resource/bootstrap: $(SRC_FILES) go.mod go.sum
	mkdir -p resource
	./run-in-docker.sh go build -tags lambda.norpc -o resource/bootstrap .

version.go template.yaml: VERSION generate.sh template.template.yaml
	./generate.sh

macro.zip: macro/app.py
	cd macro && zip -r ../macro.zip .

resource.zip: resource/bootstrap
	cd resource && zip -r ../resource.zip .

test: ## run tests
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	if command -v cfn-lint; then \
		cfn-lint --override-spec cfn-resource-specification.json example.yaml; \
	else \
		echo "cfn-lint is not found. skip it" >&2; \
	fi

clean:
	-rm -f resource.zip
	-rm -f macro.zip
	-rm -rf .build .build-sam resource
	-docker volume rm cfn-mackerel-macro-cache
