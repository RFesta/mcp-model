# {{MCP_NAME}} Makefile
# Vertikon Marketing Intelligence Platform

.PHONY: help build test run clean docker-build docker-push k8s-deploy k8s-delete

# Variables
APP_NAME := {{MCP_NAME}}
VERSION := $(shell git describe --tags --always --dirty)
DOCKER_IMAGE := vertikon/$(APP_NAME)
NAMESPACE := vertikon

# Variáveis de build
GO_VERSION = 1.21
DOCKER_BUILDKIT = 1
CGO_ENABLED = 0
GOOS = linux
GOARCH = amd64

# Arquivos e diretórios
MAIN_PATH = ./cmd/server
BINARY_NAME = main
BUILD_DIR = ./build
COVERAGE_FILE = coverage.out

# Flags de build
LDFLAGS = -w -s -extldflags "-static"
BUILD_FLAGS = -a -installsuffix cgo

# Comandos
GO = go
DOCKER = docker
KUBECTL = kubectl
HELM = helm
GOLANGCI_LINT = golangci-lint

# Cores para output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
BLUE = \033[0;34m
NC = \033[0m # No Color

.PHONY: help
help: ## Mostra este help
\t@echo "${BLUE}Makefile para $(SERVICE_NAME)${NC}"
\t@echo ""
\t@echo "${YELLOW}Comandos disponíveis:${NC}"
\t@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

####################
# Desenvolvimento
####################

.PHONY: setup
setup: ## Setup inicial do projeto
\t@echo "${BLUE}🔧 Configurando projeto $(SERVICE_NAME)...${NC}"
\t@cp .env.example .env
\t@$(GO) mod download
\t@$(GO) mod tidy
\t@echo "${GREEN}✅ Setup concluído!${NC}"

.PHONY: deps
deps: ## Instala dependências Go
\t@echo "${BLUE}📦 Instalando dependências...${NC}"
\t@$(GO) mod download
\t@$(GO) mod verify
\t@$(GO) mod tidy
\t@echo "${GREEN}✅ Dependências instaladas!${NC}"

.PHONY: dev
dev: ## Executa em modo desenvolvimento (hot reload)
\t@echo "${BLUE}🚀 Iniciando desenvolvimento...${NC}"
\t@echo "${YELLOW}Endpoints disponíveis:${NC}"
\t@echo "  - API: http://localhost:8080"
\t@echo "  - Metrics: http://localhost:9090/metrics"
\t@echo "  - Health: http://localhost:8080/public/health"
\t@echo "  - Profiling: http://localhost:6060/debug/pprof/"
\t@echo ""
\t@$(GO) run $(MAIN_PATH)

.PHONY: infra-up
infra-up: ## Sobe infraestrutura local (docker-compose)
\t@echo "${BLUE}🐳 Subindo infraestrutura local...${NC}"
\t@$(DOCKER) compose up -d postgres clickhouse redis nats minio prometheus grafana
\t@echo "${GREEN}✅ Infraestrutura disponível!${NC}"
\t@echo "${YELLOW}Serviços:${NC}"
\t@echo "  - PostgreSQL: localhost:5432"
\t@echo "  - ClickHouse: localhost:8123"
\t@echo "  - Redis: localhost:6379"
\t@echo "  - NATS: localhost:4222"
\t@echo "  - MinIO: localhost:9000 (Console: localhost:9001)"
\t@echo "  - Prometheus: localhost:9091"
\t@echo "  - Grafana: localhost:3000 (admin/admin)"

.PHONY: infra-down
infra-down: ## Para infraestrutura local
\t@echo "${BLUE}🛑 Parando infraestrutura local...${NC}"
\t@$(DOCKER) compose down
\t@echo "${GREEN}✅ Infraestrutura parada!${NC}"

####################
# Build & Test
####################

.PHONY: build
build: ## Build da aplicação
\t@echo "${BLUE}🔨 Building $(SERVICE_NAME)...${NC}"
\t@mkdir -p $(BUILD_DIR)
\t@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
\t\t$(GO) build $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" \
\t\t-o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
\t@echo "${GREEN}✅ Build concluído: $(BUILD_DIR)/$(BINARY_NAME)${NC}"

.PHONY: test
test: ## Executa todos os testes
\t@echo "${BLUE}🧪 Executando testes...${NC}"
\t@$(GO) test -v ./...
\t@echo "${GREEN}✅ Testes concluídos!${NC}"

.PHONY: coverage
coverage: ## Gera relatório de cobertura
\t@echo "${BLUE}📊 Gerando relatório de cobertura...${NC}"
\t@$(GO) test -coverprofile=$(COVERAGE_FILE) ./...
\t@$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html
\t@$(GO) tool cover -func=$(COVERAGE_FILE)
\t@echo "${GREEN}✅ Relatório gerado: coverage.html${NC}"

####################
# Quality & Security
####################

.PHONY: lint
lint: ## Executa linting
\t@echo "${BLUE}🔍 Executando linting...${NC}"
\t@$(GOLANGCI_LINT) run ./...
\t@echo "${GREEN}✅ Linting concluído!${NC}"

.PHONY: fmt
fmt: ## Formata código
\t@echo "${BLUE}🎨 Formatando código...${NC}"
\t@$(GO) fmt ./...
\t@goimports -w .
\t@echo "${GREEN}✅ Código formatado!${NC}"

####################
# Docker
####################

.PHONY: docker-build
docker-build: ## Build da imagem Docker
\t@echo "${BLUE}🐳 Building imagem Docker...${NC}"
\t@DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) $(DOCKER) build \
\t\t-t $(REGISTRY)/$(SERVICE_NAME):$(VERSION) \
\t\t-t $(REGISTRY)/$(SERVICE_NAME):latest \
\t\t--build-arg GO_VERSION=$(GO_VERSION) \
\t\t.
\t@echo "${GREEN}✅ Imagem Docker criada: $(REGISTRY)/$(SERVICE_NAME):$(VERSION)${NC}"

.PHONY: docker-run
docker-run: ## Executa container Docker
\t@echo "${BLUE}🐳 Executando container...${NC}"
\t@$(DOCKER) run --rm -p 8080:8080 -p 9090:9090 \
\t\t--env-file .env \
\t\t$(REGISTRY)/$(SERVICE_NAME):$(VERSION)

####################
# Database
####################

.PHONY: db-migrate
db-migrate: ## Executa migrations
\t@echo "${BLUE}💾 Executando migrations...${NC}"
\t@$(GO) run ./cmd/migrator up
\t@echo "${GREEN}✅ Migrations aplicadas!${NC}"

####################
# Utilities
####################

.PHONY: clean
clean: ## Limpa artefatos de build
\t@echo "${BLUE}🧹 Limpando artefatos...${NC}"
\t@rm -rf $(BUILD_DIR)
\t@rm -f $(COVERAGE_FILE) coverage.html
\t@$(GO) clean -cache -testcache -modcache
\t@echo "${GREEN}✅ Limpeza concluída!${NC}"

.PHONY: help
help: ## Mostra ajuda
\t@echo "Makefile para $(SERVICE_NAME)"

# Make sure all targets work even if files with the same name exist
.DEFAULT_GOAL := help
