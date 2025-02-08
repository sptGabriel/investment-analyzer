.PHONY: install-tools
install-tools:
	@echo "==> Installing moq"
	@go install github.com/matryer/moq@latest
	@echo "==> Installing gotest"
	@go install github.com/rakyll/gotest@latest

.PHONY: install-linters
install-linters:
	@echo "==> Installing staticcheck"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "==> Installing govulncheck"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "==> Installing GCI"
	@go install github.com/daixiang0/gci@latest

.PHONY: lint
lint:
	@echo "==> Running go vet"
	@go vet ./...
	@echo "==> Running staticcheck"
	@staticcheck ./...
	@echo "==> Running govulncheck"
	@govulncheck ./...
	@echo "==> Running gci"
	@gci write -s standard -s default -s localmodule --skip-generated .

setup: install-tools install-linters
	@go mod tidy

.PHONY: test
test:
	@echo "==> Running tests"
	@gotest -race -failfast ./...