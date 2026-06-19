# booptube

Baixe vídeos do YouTube em **mp4** ou **mp3**. Embute **yt-dlp** e **ffmpeg** — nada extra para instalar além do executável.

Disponível em **dois modos**:

| Modo | Executável | Descrição |
|------|------------|-----------|
| **CLI** | `booptube` / `booptube.exe` | Terminal interativo |
| **GUI** | `booptube-gui` / `booptube-gui.exe` | Janela gráfica neon futurista |

## Documentação

| Guia | Para quem |
|------|-----------|
| **[doc/projeto.md](doc/projeto.md)** | Visão geral completa — o que foi feito, instalação, como rodar |
| **[doc/gui.md](doc/gui.md)** | GUI — instalar, compilar e usar o `booptube-gui` |
| **[doc/usuario.md](doc/usuario.md)** | Usuário — CLI e GUI no dia a dia |
| **[doc/cli.md](doc/cli.md)** | Desenvolvedor — build, Makefile, config técnica |
| **[doc/README.md](doc/README.md)** | Índice da documentação |

## Uso rápido — GUI

```powershell
# Windows — duplo clique ou:
.\booptube-gui.exe
```

```bash
# Linux / macOS
./booptube-gui
```

1. Escolha a pasta de destino
2. Cole a URL do YouTube
3. Selecione MP4 ou MP3
4. Clique em **Baixar**

Instalação e compilação da GUI: **[doc/gui.md](doc/gui.md)**

## Uso rápido — CLI

```powershell
# Windows
.\booptube.exe
.\booptube.exe -dir "C:\Downloads"
```

```bash
# Linux / macOS
./booptube
./booptube -dir "$HOME/Downloads"
```

Loop interativo: pasta → URL → formato (`1`=mp4, `2`=mp3). Digite `q` ou `sair` para encerrar.

## Compilar (desenvolvedores)

**CLI:**

```bash
make build                    # Linux/macOS → .build/booptube
go build -o .build/booptube.exe .   # Windows
```

**GUI** (requer GCC / CGO):

```bash
make build-gui                # Linux/macOS → .build/booptube-gui
$env:CGO_ENABLED="1"; go build -tags gui -o .build/booptube-gui.exe .   # Windows
```

Antes de compilar, baixe os assets embutidos:

```powershell
.\scripts\fetch-ytdlp.ps1; .\scripts\fetch-ffmpeg.ps1   # Windows
```

```bash
make fetch-deps   # Linux/macOS
```

Detalhes em [doc/cli.md](doc/cli.md) e [doc/projeto.md](doc/projeto.md).
