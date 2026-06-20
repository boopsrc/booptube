# booptube

Baixe vídeos do YouTube em **mp4** ou **mp3**. Embute **yt-dlp** e **ffmpeg** — nada extra para instalar além do executável.

Disponível em **três modos**:

| Modo | Executável | Descrição |
|------|------------|-----------|
| **CLI** | `booptube` / `booptube.exe` | Terminal interativo |
| **GUI** | `booptube-gui` / `booptube-gui.exe` | Janela gráfica neon futurista |
| **Web** | `booptube-web` | Servidor HTTP — download no servidor, entrega via browser |

## Documentação

| Guia | Para quem |
|------|-----------|
| **[doc/projeto.md](doc/projeto.md)** | Visão geral completa — o que foi feito, instalação, como rodar |
| **[doc/gui.md](doc/gui.md)** | GUI — instalar, compilar e usar o `booptube-gui` |
| **[doc/usuario.md](doc/usuario.md)** | Usuário — CLI e GUI no dia a dia |
| **[doc/build-windows.md](doc/build-windows.md)** | Compilar no Windows (PowerShell) |
| **[doc/build-linux.md](doc/build-linux.md)** | Compilar no Linux |
| **[doc/build-macos.md](doc/build-macos.md)** | Compilar no macOS |
| **[doc/web.md](doc/web.md)** | Web — Docker, API REST, Grafana e logs |
| **[doc/cli.md](doc/cli.md)** | Desenvolvedor — referência técnica, Makefile, config |
| **[doc/installer.md](doc/installer.md)** | Releases — portable (zip) e instaladores (setup/deb/dmg) |
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

## Uso rápido — Web

```bash
cd docker
cp .env.example .env   # defina GRAFANA_ADMIN_USER e GRAFANA_ADMIN_PASSWORD
docker compose up -d --build
```

- App: http://localhost:8080
- Grafana: http://localhost:3000 (login obrigatório)

Guia completo: **[doc/web.md](doc/web.md)**

## Releases (GitHub)

Binários pré-compilados (portable + instalador) para Windows, Linux e macOS são publicados automaticamente na aba **[Releases](https://github.com/booptube/booptube/releases)** quando um commit em `main` contém **`Bump version`** no título (versão definida em [`VERSION`](VERSION)).

Guia completo: **[doc/installer.md](doc/installer.md)**

## Compilar (desenvolvedores)

| Sistema | Guia |
|---------|------|
| Windows | **[doc/build-windows.md](doc/build-windows.md)** |
| Linux | **[doc/build-linux.md](doc/build-linux.md)** |
| macOS | **[doc/build-macos.md](doc/build-macos.md)** |

Referência técnica: [doc/cli.md](doc/cli.md) · Instaladores: [doc/installer.md](doc/installer.md) · Visão geral: [doc/projeto.md](doc/projeto.md)
