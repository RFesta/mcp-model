# 📐 Blueprint de Arquitetura — Modelo MCP (Vertikon)

> Este blueprint é **estável** e deve ser seguido por todos os novos MCPs do ecossistema Vertikon.

## 1) Objetivos e Critérios de Pronto
- **Consistência** de estrutura e operações (observabilidade, logging, readiness, métricas).
- **Conectividade padronizada**: NATS JetStream (durables, consumer group), HTTP, Config.
- **Segurança básica**: sem segredos no repo, variáveis em `Secrets`/`K8s`, mínimos privilégios.
- **Operabilidade**: health/readiness, /metrics, traces, logs estruturados, pprof opcional.
- **Extensibilidade**: pontos de extensão `internal/handler`, `internal/domain`, `pkg/contracts`.
- **Deployável**: Dockerfile multi-stage + manifests K8s + Helm Chart mínimo.
- **Automação**: Makefile para build/test/lint/run + CI de verificação.

## 2) Padrões de Projeto
- **Camadas**:
  - `cmd/<service>/main.go` — composição do app
  - `internal/config` — carga e validação
  - `internal/log` — logger comum
  - `internal/otel` — tracing/metrics/pprof
  - `internal/metrics` — registries/métricas custom
  - `internal/transport/http` — roteador, handlers HTTP infra
  - `internal/nats` — client NATS/JetStream, consumidor/produtor
  - `internal/handler` — _application logic_ (assuntos NATS, core endpoints)
  - `internal/domain` — regras de negócio, entidades (quando aplicável)
  - `pkg/contracts` — schemas e DTOs estáveis (eventos/requests/responses)
  - `configs/` — YAML base
  - `deploy/` — K8s/Helm
  - `docker/` — runtime local (compose) e Dockerfile
- **Config via Viper**: Env > YAML. Prefixo `MCP_`.
- **Assuntos NATS**: `mcp.<contexto>.<serviço>.<tipo>`
  - Ex.: `mcp.modelo.example.request`, `mcp.modelo.example.reply`, `mcp.modelo.events.*`
- **Observabilidade**:
  - `GET /metrics` Prometheus
  - Tracing OTLP (por env `OTEL_EXPORTER_OTLP_ENDPOINT`)
  - `GET /debug/pprof/*` opcional, atrás de flag/env
- **Logs**: JSON com chave `service`, `env`, `trace_id`, `span_id`, `subject` (quando houver).

## 3) Pontos Estáveis (não mudam entre MCPs)
- Estrutura de pastas
- `internal/log`, `internal/otel`, `internal/metrics`, `internal/transport/http` (infra)
- Healthz/Readyz, Info
- Dockerfile, Makefile, CI base, Chart Helm básico
- Convenções de assuntos NATS e nomenclatura de variáveis

## 4) Pontos de Extensão
- `internal/handler/*` — consumidores/produtores NATS específicos
- `internal/domain/*` — regras de negócio
- `pkg/contracts/*` — esquemas/DTO de eventos/requests
- `configs/config.yaml` — subjects, timeouts e parâmetros do serviço

## 5) Fluxo de Execução (runtime)
1. Carrega config (env > YAML), inicializa logger.
2. Inicializa OTEL (tracing/metrics), registra `/metrics`.
3. Conecta ao NATS JetStream (com retry e health-check).
4. Sobe HTTP server (healthz/readyz/info/metrics/pprof).
5. Registra consumidores NATS (durables), inicia handlers.
6. Graceful shutdown: HTTP, consumidores, NATS, OTEL.

## 6) Segurança & Operação
- Sem segredos em arquivos YAML; use **K8s Secrets** ou **Vault**.
- Readiness espera conexão estável ao NATS e conclusão de warm-up dos handlers.
- HPA orientado por CPU e/ou **RPS**/Latência (via ServiceMonitor).
- Limites de recursos definidos; liveness para travamentos; pprof só em redes internas.

## 7) Versionamento e Build Info
- `internal/version` preenche **version**, **commit**, **buildTime** via `-ldflags`.
- Expostos em `/info` e nas métricas (labels).

## 8) Roadmap (sugerido, opcional por serviço)
- Outbox / Idempotência
- DLT (dead-letter) via JetStream (stream `*.DLQ`)
- Policies de retenção/replicação por criticidade

## 9) Estrutura de Pastas (referência)
```
cmd/modelo-mcp/main.go
internal/config/
internal/log/
internal/otel/
internal/metrics/
internal/transport/http/
internal/nats/
internal/handler/
internal/domain/            (opcional)
pkg/contracts/
configs/
docker/
deploy/k8s/
deploy/helm/modelo-mcp/
.github/workflows/
scripts/
```
