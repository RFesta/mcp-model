# 📖 Template Usage Guide - Modelo MCP

Este guia explica como usar o template `modelo-mcp` para criar novos MCPs rapidamente no ecossistema Vertikon.

## 🚀 Como Usar Este Template

### 1. Copiar Template
```bash
# Copiar para novo MCP
cp -r E:\vertikon\modelo-mcp E:\vertikon\marketing\criacao\mcp-content-generator

# Navegar para o novo MCP
cd E:\vertikon\marketing\criacao\mcp-content-generator
```

### 2. Substituir Placeholders
Use find & replace no seu editor para substituir estas variáveis:

#### Variáveis Principais
- `{{MCP_NAME}}` → `mcp-content-generator`
- `{{MCP_MODULE_NAME}}` → `mcp-content-generator`
- `{{MCP_DESCRIPTION}}` → `AI-powered content generation service`
- `{{DATABASE_NAME}}` → `content_generator`
- `{{CLICKHOUSE_DATABASE}}` → `content_analytics`

#### Variáveis de Serviço
- `{{AI_SERVICE_1}}` → `ContentGenerator`
- `{{AI_SERVICE_2}}` → `QualityAnalyzer`
- `{{AI_SERVICE_3}}` → `SEOOptimizer`
- `{{AI_SERVICE_4}}` → `PerformanceTracker`

#### Variáveis de Negócio
- `{{CORE_SERVICE}}` → `ContentManagement`
- `{{MAIN_ENTITY}}` → `Content`
- `{{MAIN_ENDPOINT}}` → `content`

#### Configurações Específicas
- `{{SERVICE_CONFIG_NAME}}` → `ContentConfig`
- `{{SERVICE_CONFIG_KEY}}` → `content_generator`

### 3. Personalizar Funcionalidades

#### 3.1 Atualizar main.go
```go
// Substitua as seções específicas do serviço:
// - Inicialização de clientes específicos
// - AI services personalizados
// - Endpoints específicos do domínio
// - Background tasks específicos
```

#### 3.2 Criar Models Específicos
```go
// internal/models/content.go
type Content struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title"`
    Body        string    `json:"body"`
    Type        string    `json:"type"`
    Status      string    `json:"status"`
    TenantID    string    `json:"tenant_id" gorm:"index"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### 3.3 Implementar Services
```go
// internal/services/content_service.go
type ContentService struct {
    db          *gorm.DB
    clickhouse  clickhouse.Conn
    aiGenerator AIContentGeneratorService
    config      *config.Config
}
```

#### 3.4 Criar Handlers
```go
// internal/handlers/content_handler.go
func (h *ContentHandler) CreateContent(c *gin.Context) {
    // Implementar lógica específica
}
```

### 4. Configurar Dependências

#### 4.1 go.mod
```bash
# Adicionar dependências específicas se necessário
go get github.com/specific-package
go mod tidy
```

#### 4.2 config.yaml
```yaml
content_generator:
  max_words: 1000
  quality_threshold: 0.8
  seo_optimization: true
  supported_formats:
    - "blog_post"
    - "social_media"
    - "email"
```

### 5. Implementar Métricas Específicas

#### 5.1 Adicionar em pkg/metrics/metrics.go
```go
var (
    ContentGenerated = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mcp_content_generator_content_generated_total",
            Help: "Total content generated",
        },
        []string{"type", "status", "tenant_id"},
    )
    
    GenerationDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "mcp_content_generator_generation_duration_seconds",
            Help: "Content generation duration",
        },
        []string{"type"},
    )
)
```

## 🔧 Customizações Comuns

### Para MCPs de Analytics
```go
// Adicionar mais integrações ClickHouse
// Implementar dashboards específicos
// Criar agregações personalizadas
```

### Para MCPs de Automação
```go
// Adicionar workers para background processing
// Implementar workflows específicos
// Integrar com sistemas externos
```

### Para MCPs de AI/ML
```go
// Configurar providers específicos de AI
// Implementar pipelines de ML
// Adicionar modelo training/inference
```

## 📁 Arquivos que Sempre Precisam Personalização

### Obrigatórios
- ✅ `main.go` - Lógica principal e endpoints
- ✅ `internal/models/` - Models específicos do domínio
- ✅ `internal/services/` - Business logic
- ✅ `internal/handlers/` - HTTP handlers
- ✅ `configs/config.yaml` - Configurações específicas
- ✅ `README.md` - Documentação do MCP

### Opcionais (dependendo do MCP)
- `pkg/metrics/metrics.go` - Métricas específicas
- `migrations/` - Schema específico
- `k8s/` - Configurações específicas de deploy
- `docker-compose.yml` - Serviços específicos

## 🎯 Exemplo Completo: MCP Content Generator

### Estrutura Final
```
mcp-content-generator/
├── main.go                                    # ✅ Personalizado
├── internal/
│   ├── models/content.go                      # ✅ Específico
│   ├── services/content_service.go            # ✅ Específico
│   ├── handlers/content_handler.go            # ✅ Específico
│   ├── config/config.go                       # ✅ Template usado
│   ├── database/database.go                   # ✅ Template usado
│   └── middleware/middleware.go               # ✅ Template usado
├── pkg/
│   ├── logger/logger.go                       # ✅ Template usado
│   └── metrics/metrics.go                     # ✅ Personalizado
├── configs/config.yaml                        # ✅ Personalizado
├── docker-compose.yml                         # ✅ Template usado
├── Dockerfile                                 # ✅ Template usado
├── go.mod                                     # ✅ Personalizado
└── README.md                                  # ✅ Personalizado
```

### Comandos de Desenvolvimento
```bash
# Setup inicial
make deps
make infra-up

# Desenvolvimento
make dev

# Testes
make test
make coverage

# Build
make build
make docker-build

# Deploy
make k8s-deploy
```

## 💡 Dicas de Produtividade

1. **Use Search & Replace Global**: Para substituir todas as variáveis de uma vez
2. **Copie Lógica Similar**: Veja MCPs existentes similares para copiar patterns
3. **Teste Incremental**: Teste cada funcionalidade conforme implementa
4. **Documente APIs**: Mantenha README.md atualizado com novos endpoints
5. **Métricas Desde Início**: Adicione métricas de negócio desde a primeira versão

---

**🎯 Com este template, você pode criar novos MCPs em 15-30 minutos!**