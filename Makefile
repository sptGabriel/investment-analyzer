BIN_FOLDER=./bin
YQ_BIN := $(shell command -v yq 2>/dev/null)
YQ_GOPATH_BIN := $(shell go env GOPATH)/bin/yq
DOCKER_COMPOSE_CMD=docker compose

.PHONY: ensure-yq
ensure-yq:
ifeq ($(YQ_BIN),)
	@echo "==> yq não encontrado. Instalando via go..."
	@go install github.com/mikefarah/yq/v4@latest
	@echo "==> yq instalado em $(YQ_GOPATH_BIN)."
else
	@echo "==> yq já está instalado em $(YQ_BIN)."
endif

.PHONY: install-tools
install-tools:
	@echo "==> Installing moq"
	@go install github.com/matryer/moq@latest
	@echo "==> Installing gotest"
	@go install github.com/rakyll/gotest@latest
	@echo "==> Installing swaggo"
	@go install github.com/swaggo/swag/cmd/swag@v1.8.12

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

.PHONY: test
test:
	@echo "==> Running tests"
	@gotest -race -failfast ./...

.PHONY: generate-analyzer-swagger
generate-swagger:
	@echo "==> Formating swagger doc"
	@swag  fmt  -d ./app/analyzer/api/v1 -g api.go
	@echo "==> Generating swagger doc"
	@swag  init --parseDependency --parseInternal -d ./app/analyzer/api/v1 -g api.go -o docs/swagger

.PHONY: generate-mocks
generate-mocks:
	@echo "==> Generating mocks"
	@go generate ./...

.PHONY: generate
generate: generate-sqlc generate-swagger generate-mocks
	@echo "==> Generate done!"


.PHONY: build
build: ensure-yq
	@mkdir -p $(BIN_FOLDER)
	@appname=$$($(or $(YQ_BIN),$(YQ_GOPATH_BIN)) '.application' buildfile.yaml); \
	yqMainPaths=$$($(or $(YQ_BIN),$(YQ_GOPATH_BIN)) '.binaries[] | .main' buildfile.yaml); \
	echo "Aplicação: $$appname"; \
	for codefolder in $$yqMainPaths; do \
		binary=$$appname-$$($(or $(YQ_BIN),$(YQ_GOPATH_BIN)) '.binaries[] | select(.main == "'$$codefolder'") | .name' buildfile.yaml); \
		echo "==> Build $$binary from folder $$codefolder"; \
		GOOS=linux CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-s -w" -o $(BIN_FOLDER)/$$binary $$codefolder; \
	done
	@echo "==> Build done!"


.PHONY: deps-check-go-version
deps-check-go-version: deps-check-go-install buildfile.yaml
	@required_version=$$(grep 'go-version' buildfile.yaml | awk '{print $$2}' | tr -d '"'); \
	current_version=$$(go version | awk '{print $$3}' | sed 's/go//'); \
	if [ "$$(printf '%s\n' "$$required_version" "$$current_version" | sort -V | head -n1)" = "$$required_version" ]; then \
		echo "\t\033[1;32mGo version is OK (+$$required_version).\033[0m"; \
	else \
		echo "\033[1;31mGo version is less than $$required_version. Please update Go.\033[0m"; \
		exit 1; \
	fi

.PHONY: deps-check-go-install
deps-check-go-install:
	@command -v go >/dev/null 2>&1 || { echo >&2 "\033[1;31mGo is not installed. Please install Go to proceed.\033[0m"; exit 1; }
	@echo "\t\033[1;32mGo is installed.\033[0m"

.PHONY: deps-check-docker
deps-check-docker:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "\033[1;31mDocker is not installed. Please install Docker to proceed.\033[0m"; exit 1; }
	@echo "\t\033[1;32mDocker is installed.\033[0m"

.PHONY: dependences
dependences: deps-check-go-version deps-check-docker
	@command -v yq >/dev/null 2>&1 || { echo >&2 "\033[1;31myq is not installed. Please install yq to proceed.\033[0m"; exit 1; }
	@echo "\033[1;32mAll dependencies are installed.\033[0m"

.PHONY: setup
setup: .env dependences install-tools install-linters
	@echo "=> Setting up the project"
	@echo "===> Pulling docker image: xanders/make-help:latest \033[1;33m(warn:latest-tag)\033[0m"
	@docker pull xanders/make-help:latest
	@echo "===> Pulling docker image: compose images"
	@$(DOCKER_COMPOSE_CMD) pull

.PHONY: docker-up
docker-up:
	@docker ps -a | awk '{print $$1}' | grep -v CONTAINER | xargs -I {} docker stop {}
	@docker ps -a | awk '{print $$1}' | grep -v CONTAINER | xargs -I {} docker rm {}
	@echo "===> Building and starting containers"
	@$(DOCKER_COMPOSE_CMD) up -d
	@echo "\033[1;32m running!\033[0m"