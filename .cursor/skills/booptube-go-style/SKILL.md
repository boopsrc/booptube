---
name: booptube-go-style
description: >-
  Enforces flat Go architecture and minimal conventions for the booptube project:
  context-based folders, no single-impl interfaces, consolidated files, inline
  errors, and no obvious comments. Use when writing, refactoring, reviewing, or
  scaffolding Go code in the booptube repository.
---

# Booptube Go Style

## Estrutura do Projeto

A raiz contém os contextos de negócio e configuração essencial. Entry points ficam em `cmd/`. Evite `pkg/` ou `internal/` desnecessários.

```
booptube/
├── cmd/
│   ├── cli/main.go      # booptube (terminal)
│   └── gui/main.go      # booptube-gui (Fyne)
├── downloader/          # Contexto 1: Engine de download
│   ├── client.go
│   └── client_test.go
├── video/               # Contexto 2: Domínio/Regras de negócio do vídeo
│   └── model.go
├── ui/                  # Contexto 3: Interface
│   ├── terminal.go      # CLI (ui.Run)
│   └── gui/             # GUI Fyne (gui.Run)
├── config/              # Contexto 4: Configurações globais
│   └── config.go
└── go.mod
```

### Onde colocar código novo

| Contexto | Pasta | Exemplos |
|----------|-------|----------|
| Download, HTTP, streams | `downloader/` | `client.go` |
| Modelos e regras de vídeo | `video/` | `model.go` |
| CLI, prompts, output | `ui/` | `terminal.go` |
| GUI Fyne | `ui/gui/` | `gui.go`, `theme.go` |
| Env, flags, defaults | `config/` | `config.go` |
| Wiring e bootstrap | `cmd/cli/`, `cmd/gui/` | `main.go` |

**Não criar:** `pkg/`, `internal/`, `api/`, `services/`, `handlers/`, `controllers/`, `repositories/`.

## Regras de Código

### Sem Interfaces de Implementação Única

Não crie `type Service interface` se houver apenas uma struct que a implementa. Use a struct diretamente.

```go
// ❌ Evitar
type Downloader interface { Download(url string) error }
type Client struct{}
func (c *Client) Download(url string) error { ... }

// ✅ Preferir
type Client struct{}
func (c *Client) Download(url string) error { ... }
```

Interfaces só quando há **duas ou mais** implementações reais, ou quando o pacote precisa de mock explícito em testes.

### Arquivos Consolidados por Contexto

Em vez de criar `router.go`, `handler.go` e `controller.go`, aglutine funções correlatas em um único arquivo (ex: `client.go`). Menos caminhos de arquivo = menos tokens de contexto.

- Funções do mesmo contexto e responsabilidade → mesmo arquivo
- Novo arquivo só quando o arquivo atual ficar difícil de navegar (~300+ linhas) ou quando o subdomínio for claramente distinto

### Tratamento de Erros Inline Simplificado

Evite criar tipos de erro customizados complexos se `fmt.Errorf` ou `errors.New` resolverem.

```go
// ❌ Evitar (para erros simples)
type ErrNotFound struct { URL string }
func (e ErrNotFound) Error() string { ... }

// ✅ Preferir
return fmt.Errorf("video not found: %s", url)
```

Tipos de erro customizados apenas quando o caller precisa fazer `errors.Is` / `errors.As` com semântica de negócio.

### Zero Comentários Óbvios

Remova comentários como `// NewClient cria um novo cliente`. O código deve ser autoexplicativo. Use comentários apenas para hacks/regras de negócio obscuras.

```go
// ❌ Evitar
// Download baixa o vídeo da URL informada.
func (c *Client) Download(url string) error

// ✅ Preferir (sem comentário)
func (c *Client) Download(url string) error

// ✅ OK — regra de negócio não óbvia
// YouTube retorna 403 se o User-Agent não imitar um browser real.
```

## Checklist antes de entregar

- [ ] Código está no contexto correto (`downloader/`, `video/`, `ui/`, `config/`)
- [ ] Nenhuma pasta `pkg/`, `internal/` ou camada extra desnecessária
- [ ] Sem interface com implementação única
- [ ] Funções correlatas no mesmo arquivo, não espalhadas em múltiplos arquivos finos
- [ ] Erros com `fmt.Errorf` / `errors.New`, sem tipos customizados triviais
- [ ] Sem comentários que apenas repetem o nome da função ou tipo
- [ ] `cmd/*/main.go` apenas faz wiring; lógica de negócio fica nos contextos
