# 🚀 {{MCP_NAME}} - Marketing Intelligence Platform

> **{{MCP_DESCRIPTION}}** with AI-powered automation and real-time analytics.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Coverage](https://img.shields.io/badge/coverage-85%25-yellowgreen.svg)]()

---

## 🎯 CARACTERÍSTICAS

### ✅ Stack Tecnológico Completo
- **Backend**: Go 1.21+ com Gin/Fiber framework
- **Database**: PostgreSQL 16+ com RLS + ClickHouse para analytics  
- **Cache**: Redis Cluster com distribuição automática
- **Message Queue**: NATS JetStream para processamento assíncrono
- **Storage**: MinIO/S3 compatível para assets
- **Monitoring**: Prometheus + Grafana + custom metrics
- **Security**: JWT Auth + rate limiting + encryption

### ✅ DevOps & Infraestrutura
- **Containerização**: Docker multi-stage builds otimizados
- **Orquestração**: Kubernetes + Helm charts
- **CI/CD**: GitHub Actions com pipeline completo
- **Observabilidade**: Logs estruturados + tracing + metrics
- **Quality Gates**: Testes automatizados + security scan + lint

### ✅ Enterprise Features
- **Multi-tenancy**: Isolamento completo por tenant com RLS
- **High Availability**: Auto-scaling + health checks + graceful shutdown
- **Security**: OWASP compliance + audit trail + data encryption
- **Performance**: Connection pooling + caching + optimizations

---

## 🚀 QUICK START

### 1. Clone do Template
```bash
# Método 1: Via GitHub CLI (recomendado)
gh repo create meu-novo-mcp --template vertikon/modelo-mcp --private
cd meu-novo-mcp

# Método 2: Clone manual
git clone https://github.com/vertikon/modelo-mcp.git meu-novo-mcp
cd meu-novo-mcp
rm -rf .git
git init
```

### 2. Setup Automatizado
```bash
# Script PowerShell (Windows)
.\scripts\setup.ps1 -ServiceName "mcp-meu-servico" -Port 8080

# Script Bash (Linux/Mac)
./scripts/setup.sh mcp-meu-servico 8080

# Setup manual
make setup SERVICE_NAME=mcp-meu-servico PORT=8080
```

### 3. Desenvolvimento Local
```bash
# Instalar dependências
make deps

# Subir infraestrutura local (Docker Compose)
make infra-up

# Executar em modo desenvolvimento
make dev

# Executar testes
make test

# Build da aplicação
make build
```

### 4. Deploy para Kubernetes
```bash
# Deploy completo
make k8s-deploy NAMESPACE=vertikon-prod

# Deploy apenas da aplicação
helm upgrade --install mcp-meu-servico ./helm/mcp-template

# Monitoramento
kubectl logs -f deployment/mcp-meu-servico -n vertikon-prod
```

---

## 📁 ESTRUTURA DO PROJETO

```
modelo-mcp/
├── cmd/
│   ├── server/              # Entrypoint principal
│   ├── migrator/            # Migrations database
│   └── worker/              # Background workers
├── internal/
│   ├── api/                 # HTTP handlers e routes
│   ├── config/              # Configuration management
│   ├── database/            # Database connections
│   ├── middleware/          # HTTP middlewares
│   ├── models/              # Database models (GORM)
│   ├── services/            # Business logic
│   └── workers/             # Background jobs
├── pkg/
│   ├── logger/              # Structured logging
│   ├── metrics/             # Prometheus metrics
│   ├── storage/             # File storage (S3/MinIO)
│   └── utils/               # Utilities
├── migrations/              # SQL migrations
├── helm/                    # Helm charts
├── k8s/                     # Kubernetes manifests
├── scripts/                 # Setup e utility scripts
├── docker/                  # Dockerfiles
├── docs/                    # Documentação
├── tests/                   # Testes (unit, integration, e2e)
├── docker-compose.yml       # Desenvolvimento local
├── Makefile                 # Build automation
├── Dockerfile               # Production image
└── .github/workflows/       # CI/CD GitHub Actions
```

---

## ⚙️ CONFIGURAÇÃO

### Variáveis de Ambiente Essenciais
```bash
# Service Configuration
SERVICE_NAME=mcp-template
HTTP_PORT=8080
METRICS_PORT=9090
ENVIRONMENT=development

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/template
CLICKHOUSE_URL=http://localhost:8123/analytics
REDIS_URL=redis://localhost:6379/0

# Message Queue
NATS_URL=nats://localhost:4222
NATS_CLUSTER_ID=vertikon-cluster

# Security
JWT_SECRET=your-super-secret-key-32-chars
ENCRYPTION_KEY=32-character-encryption-key-here

# Storage
STORAGE_PROVIDER=MINIO
STORAGE_ENDPOINT=localhost:9000
STORAGE_ACCESS_KEY_ID=minioadmin
STORAGE_SECRET_ACCESS_KEY=minioadmin
STORAGE_BUCKET=mcp-assets

# Rate Limiting
RATE_LIMIT_RPM=1000
RATE_LIMIT_BURST=100

# Monitoring
ENABLE_METRICS=true
ENABLE_TRACING=false
LOG_LEVEL=info
LOG_FORMAT=json
```

### Configuração Multi-Ambiente
```bash
# Desenvolvimento
cp .env.example .env.development

# Staging  
cp .env.example .env.staging

# Produção
cp .env.example .env.production
```

---

## 🛠️ COMANDOS MAKE

```bash
# Desenvolvimento
make dev              # Executar em modo desenvolvimento
make deps             # Instalar dependências
make infra-up         # Subir infraestrutura local
make infra-down       # Parar infraestrutura local

# Build & Test
make build            # Build da aplicação
make test             # Executar todos os testes
make test-unit        # Testes unitários
make test-integration # Testes de integração
make coverage         # Relatório de cobertura

# Quality & Security
make lint             # Linting (golangci-lint)
make fmt              # Format code (gofmt + goimports)
make security-scan    # Security scanning
make dependency-check # Verificar vulnerabilidades

# Docker
make docker-build     # Build da imagem Docker
make docker-push      # Push para registry
make docker-run       # Executar container

# Kubernetes
make k8s-deploy       # Deploy completo
make k8s-undeploy     # Remove deployment
make k8s-logs         # Visualizar logs
make k8s-describe     # Descrever recursos

# Database
make db-migrate       # Executar migrations
make db-rollback      # Rollback última migration
make db-reset         # Reset database
make db-seed          # Popular com dados de teste

# Utilities
make clean            # Limpar artifacts
make help             # Mostrar todos os comandos
```

---

**🎯 Template Modelo MCP - Acelere o desenvolvimento de microserviços enterprise-grade!**

*Desenvolvido com ❤️ pelo time Vertikon*
