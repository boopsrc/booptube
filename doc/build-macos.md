# Compilar no macOS

Guia passo a passo para build da **CLI** (`booptube`) e **GUI** (`booptube-gui`) no macOS.

Fluxo recomendado: **Makefile** + scripts `.sh`. A GUI usa **CGO** com Clang do Xcode Command Line Tools (não precisa de MinGW).

Outros guias: [build-windows.md](build-windows.md) · [build-linux.md](build-linux.md)

---

## Pré-requisitos

| Requisito | Versão | CLI | GUI |
|-----------|--------|-----|-----|
| Go | 1.22+ | Sim | Sim |
| Xcode Command Line Tools | — | Não | Sim (CGO / Fyne) |
| make, bash, curl, unzip | — | Sim | Sim |

Instale Go em [go.dev/dl](https://go.dev/dl/) ou via Homebrew: `brew install go`.

### Xcode Command Line Tools (obrigatório para GUI)

```bash
xcode-select --install
```

Confirme:

```bash
clang --version
```

---

## Apple Silicon vs Intel

Os scripts fetch baixam binários embutidos para **darwin-arm64** (Apple Silicon). Compile na mesma arquitetura da máquina — o Go gera o executável nativo automaticamente.

---

## 1. Baixar yt-dlp e ffmpeg

```bash
cd /caminho/para/booptube
chmod +x scripts/*.sh
make fetch-deps
```

Ou manualmente:

```bash
./scripts/fetch-ytdlp.sh
./scripts/fetch-ffmpeg.sh
```

Os binários ficam em `assets/ytdlp/` e `assets/ffmpeg/`.

---

## 2. Compilar a CLI

**Com Makefile (recomendado):**

```bash
make build
```

**Manual:**

```bash
mkdir -p .build
go build -o .build/booptube ./cmd/cli
```

Saída: `.build/booptube` (~200 MB)

> `make build` já executa `fetch-deps` automaticamente.

---

## 3. Compilar a GUI

**Com Makefile (recomendado):**

```bash
make build-gui
```

**Manual:**

```bash
CGO_ENABLED=1 go build -o .build/booptube-gui ./cmd/gui
```

Saída: `.build/booptube-gui` (~250 MB)

A primeira compilação da GUI pode demorar (Fyne + CGO).

---

## 4. Limpar build

```bash
make clean
```

Ou:

```bash
rm -rf .build
```

---

## 5. Rodar

```bash
./.build/booptube
./.build/booptube -dir "$HOME/Downloads/booptube"
./.build/booptube-gui
```

### Instalar no PATH (opcional)

```bash
sudo install -m 755 .build/booptube .build/booptube-gui /usr/local/bin/
```

### Gatekeeper (app não abre)

Se o macOS bloquear o binário compilado localmente:

```bash
xattr -cr .build/booptube-gui
```

Ou: **Ajustes do Sistema → Privacidade e Segurança → Abrir mesmo assim**.

---

## Problemas comuns

| Erro | Solução |
|------|---------|
| `yt-dlp embutido ausente` | `make fetch-deps` |
| `xcrun: error: invalid active developer path` | `xcode-select --install` |
| `clang: command not found` | Instale Xcode CLT |
| GUI não abre (Gatekeeper) | `xattr -cr .build/booptube-gui` |

---

## Referência Makefile

| Comando | Descrição |
|---------|-----------|
| `make fetch-deps` | Baixa yt-dlp e ffmpeg |
| `make build` | fetch + compila `.build/booptube` |
| `make build-gui` | fetch + compila `.build/booptube-gui` (CGO) |
| `make clean` | Remove `.build/` |

Variáveis opcionais: `YTDLP_VERSION`, `FFMPEG_VERSION` (ver [cli.md](cli.md)).

---

## Ver também

- [gui.md](gui.md) — usar a interface gráfica
- [cli.md](cli.md) — referência técnica
- [projeto.md](projeto.md) — visão geral do projeto
