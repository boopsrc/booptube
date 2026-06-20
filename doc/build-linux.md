# Compilar no Linux

Guia passo a passo para build da **CLI** (`booptube`) e **GUI** (`booptube-gui`) no Linux.

Fluxo recomendado: **Makefile** + scripts `.sh`.

Outros guias: [build-windows.md](build-windows.md) · [build-macos.md](build-macos.md)

---

## Pré-requisitos

| Requisito | Versão | CLI | GUI |
|-----------|--------|-----|-----|
| Go | 1.22+ | Sim | Sim |
| make, bash, curl, unzip, tar | — | Sim | Sim |
| GCC, pkg-config, deps OpenGL/X11 | — | Não | Sim (CGO / Fyne) |

**yt-dlp** e **ffmpeg** vêm embutidos — não instale no sistema para compilar.

### Dependências da GUI (Fyne)

**Debian / Ubuntu:**

```bash
sudo apt update
sudo apt install -y build-essential pkg-config libgl1-mesa-dev xorg-dev
```

**Fedora / RHEL:**

```bash
sudo dnf install -y gcc pkg-config libX11-devel libXcursor-devel libXrandr-devel \
  libXinerama-devel libXi-devel libXxf86vm-devel mesa-libGL-devel libxcursor-devel
```

**Arch Linux:**

```bash
sudo pacman -S base-devel pkg-config mesa libxcursor libxrandr libxinerama libxi libxxf86vm
```

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

Saída: `.build/booptube` (~185–195 MB)

> `make build` aplica `-trimpath -ldflags "-s -w"` e injeta versão de [`VERSION`](../VERSION). A maior parte do tamanho vem dos binários embutidos (ffmpeg + yt-dlp).

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

Saída: `.build/booptube-gui` (~235–245 MB)

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

---

## Problemas comuns

| Erro | Solução |
|------|---------|
| `yt-dlp embutido ausente` | `make fetch-deps` ou rode os scripts fetch |
| `pkg-config not found` | Instale `pkg-config` (pacote do sistema) |
| Erros de OpenGL / X11 ao compilar GUI | Instale `libgl1-mesa-dev xorg-dev` (Debian) ou equivalente |
| `permission denied` nos scripts | `chmod +x scripts/*.sh` |

---

## Referência Makefile

| Comando | Descrição |
|---------|-----------|
| `make fetch-deps` | Baixa yt-dlp e ffmpeg |
| `make build` | fetch + compila `.build/booptube` |
| `make build-gui` | fetch + compila `.build/booptube-gui` (CGO) |
| `make clean` | Remove `.build/` |

Variáveis opcionais: `YTDLP_VERSION`, `FFMPEG_VERSION` (ver [cli.md](cli.md)).

Instaladores e zip portable: **[installer.md](installer.md)** (`make package-portable-linux`, `make package-linux`).

---

## Ver também

- [gui.md](gui.md) — usar a interface gráfica
- [cli.md](cli.md) — referência técnica
- [projeto.md](projeto.md) — visão geral do projeto
