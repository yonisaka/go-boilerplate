GOPATH=$(shell go env GOPATH)

GOLANGCI_LINT_VERSION=v1.53

install: install-gomock install-linter

install-gomock:
	@echo "\n>>> Install gomock\n"
	go install github.com/golang/mock/mockgen

install-linter:
	@echo "\n>>> Install GolangCI-Lint"
	@echo ">>> https://github.com/golangci/golangci-lint/releases \n"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/${GOLANGCI_LINT_VERSION}/install.sh | \
	sh -s -- -b ${GOPATH}/bin ${GOLANGCI_LINT_VERSION}

lint:
	@echo "\n>>> Run GolangCI-Lint\n"
	/bin/bash ./scripts/lint.sh

test:
	mkdir -p .coverage/html
	go test -v -race -cover -coverprofile=.coverage/internal.coverage ./internal/... && \
	cat .coverage/internal.coverage | grep -v "_mock.go\|_mockgen.go" > .coverage/internal.mockless.coverage && \
	go tool cover -html=.coverage/internal.mockless.coverage -o .coverage/html/internal.coverage.html;

mock:
	@echo "\n>>> Run Generate Mock\n"
	go generate ./...

http:
	@go run main.go http:start