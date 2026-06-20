# booptube — modo Web

Servidor HTTP em Go que baixa vídeos do YouTube **no servidor** e entrega o arquivo ao usuário via browser. Suporta múltiplos usuários simultâneos; arquivos são removidos automaticamente **10 minutos** após ficarem prontos.

| Guia relacionado | Conteúdo |
|------------------|----------|
| [projeto.md](projeto.md) | Visão geral do projeto |
| [cli.md](cli.md) | Referência técnica CLI/GUI |

---

## O que é o booptube-web

| Componente | Descrição |
|------------|-----------|
| **Página web** | UI neon em português — URL, formato MP4/MP3, progresso |
| **API REST** | Enfileira downloads, consulta status, serve arquivo |
| **Job de limpeza** | Remove arquivos após TTL configurável (padrão 10 min) |
| **Observabilidade** | Logs JSON → Loki → Grafana |

O download usa o mesmo motor [`downloader.Client`](../downloader/client.go) da CLI e GUI (yt-dlp + ffmpeg embutidos).

---

## Subir com Docker (recomendado)

### Pré-requisitos

- Docker e Docker Compose v2
- Portas livres: `8080` (app), `3000` (Grafana)

### Passos

```bash
cd docker
cp .env.example .env
```

Edite `.env` e defina credenciais **obrigatórias** do Grafana:

```env
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=sua-senha-forte-aqui
```

Suba o stack:

```bash
docker compose up -d --build
```

### Endpoints

| URL | Descrição |
|-----|-----------|
| http://localhost:8080 | Interface web |
| http://localhost:8080/health | Health check |
| http://localhost:3000 | Grafana (login obrigatório) |

### Parar

```bash
docker compose down
```

Para remover volumes (downloads e dados do Loki/Grafana):

```bash
docker compose down -v
```

---

## Variáveis de ambiente

### App (`booptube-web`)

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `BOOPTUBE_ADDR` | `:8080` | Endereço HTTP |
| `BOOPTUBE_DOWNLOAD_DIR` | `/data/downloads` | Pasta de jobs no servidor |
| `BOOPTUBE_MAX_CONCURRENT` | `4` | Downloads simultâneos |
| `BOOPTUBE_FILE_TTL` | `10m` | Tempo até expirar após pronto |
| `BOOPTUBE_DOWNLOAD_TIMEOUT` | `30m` | Timeout por download |
| `GOMEMLIMIT` | `1800MiB` | Limite de memória do GC Go |

### Grafana

| Variável | Obrigatória | Descrição |
|----------|-------------|-----------|
| `GRAFANA_ADMIN_USER` | Sim | Usuário admin |
| `GRAFANA_ADMIN_PASSWORD` | Sim | Senha admin |
| `GRAFANA_PORT` | Não (3000) | Porta no host |
| `GRAFANA_ROOT_URL` | Não | URL pública do Grafana |

---

## API REST

### `POST /api/downloads`

Inicia um download.

**Body:**

```json
{
  "url": "https://www.youtube.com/watch?v=...",
  "format": "mp4"
}
```

`format`: `"mp4"` ou `"mp3"`.

**Resposta (202):**

```json
{
  "id": "a1b2c3...",
  "status": "queued"
}
```

### `GET /api/downloads/{id}`

Consulta status do job.

**Resposta (200):**

```json
{
  "id": "a1b2c3...",
  "url": "https://...",
  "format": "mp4",
  "status": "downloading",
  "progress": 45.2,
  "log": ["[download] ..."],
  "created_at": "2026-06-19T12:00:00Z",
  "ready_at": "",
  "expires_at": "",
  "download_url": ""
}
```

Status possíveis: `queued`, `downloading`, `ready`, `failed`, `expired`.

Quando `ready`:

```json
{
  "status": "ready",
  "progress": 100,
  "filename": "Titulo do video.mp4",
  "download_url": "/api/downloads/a1b2c3.../file",
  "expires_at": "2026-06-19T12:10:00Z"
}
```

### `GET /api/downloads/{id}/file`

Baixa o arquivo gerado.

- `200` — arquivo anexado
- `409` — ainda não disponível
- `410` — expirado ou removido

---

## Logs e Grafana

O app emite logs estruturados JSON no stdout (`slog`). O **Promtail** coleta logs dos containers Docker e envia ao **Loki**. O **Grafana** já vem provisionado com:

- Datasource Loki
- Dashboard **Booptube Web** (downloads, falhas, requisições HTTP)

### Queries Loki úteis

Todos os logs do app:

```logql
{service="booptube-web"}
```

Filtrar por job:

```logql
{service="booptube-web"} | json | job_id="SEU_JOB_ID"
```

Downloads concluídos:

```logql
{service="booptube-web"} | json | msg="download_ready"
```

Falhas:

```logql
{service="booptube-web"} | json | msg="download_failed"
```

---

## Compilar localmente (sem Docker)

```bash
make build-web
```

Requer bash (Git Bash ou WSL no Windows) para `fetch-deps`.

Executar:

```bash
# Linux / macOS
BOOPTUBE_DOWNLOAD_DIR=./tmp-downloads ./.build/booptube-web

# Windows
$env:BOOPTUBE_DOWNLOAD_DIR=".\tmp-downloads"
.\.build\booptube-web.exe
```

Abra http://localhost:8080

---

## Arquitetura

```text
Browser
  → POST /api/downloads (enfileira)
  → GET  /api/downloads/{id} (polling)
  → GET  /api/downloads/{id}/file (download)

booptube-web
  → fila FIFO + semáforo (max concurrent)
  → downloader.DownloadTo(jobDir)
  → yt-dlp + ffmpeg
  → TTL 10min → remove arquivo
```

Cada job usa pasta isolada: `{BOOPTUBE_DOWNLOAD_DIR}/{jobID}/`.

---

## Dependências de runtime (Docker/Linux)

O container instala as bibliotecas de sistema necessárias para executar as ferramentas embutidas:

| Ferramenta | Tipo | Observação |
|------------|------|------------|
| **yt-dlp** | Binário standalone (`yt-dlp_linux`) | Não precisa de Python instalado |
| **ffmpeg / ffprobe** | Build estático (johnvansickle) | Sem dependências dinâmicas |
| **Libs do sistema** | `libstdc++6`, `libssl3`, `libz1`, etc. | HTTPS, TLS e runtime do yt-dlp |

Na inicialização, o servidor executa `VerifyTools()` — extrai os binários embutidos e valida `yt-dlp --version`, `ffmpeg -version` e `ffprobe -version`. Se alguma ferramenta falhar, o container encerra com log `tools_verify_error`.

No build da imagem Docker, as ferramentas também são validadas após `fetch-deps`.

---

## Produção

- Use senha forte no Grafana e HTTPS via proxy reverso (Caddy, Traefik, nginx)
- Ajuste `GOMEMLIMIT` conforme RAM do container
- Monitore disco do volume `booptube-downloads`
- A app web é pública por padrão (sem autenticação de usuários)
