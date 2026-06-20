# booptube — documentação do projeto

Visão geral de **tudo que foi implementado** até o momento: o que o booptube faz, como está organizado, como instalar e como rodar a **CLI** e a **GUI**.

| Guia relacionado | Conteúdo |
|------------------|----------|
| [usuario.md](usuario.md) | Uso do dia a dia (CLI e GUI) |
| [build-windows.md](build-windows.md) | Compilar no Windows |
| [build-linux.md](build-linux.md) | Compilar no Linux |
| [build-macos.md](build-macos.md) | Compilar no macOS |
| [cli.md](cli.md) | Referência técnica para desenvolvedores |
| [README.md](README.md) | Índice da pasta `doc/` |

---

## O que é o booptube

O **booptube** é um aplicativo para baixar vídeos do YouTube em **MP4** (vídeo + áudio) ou **MP3** (só áudio). Ele embute **yt-dlp** e **ffmpeg** dentro do executável — não é necessário instalar essas ferramentas no sistema.

Existem **dois programas**, compilados a partir do mesmo código:

| Programa | Interface | Executável |
|----------|-----------|------------|
| **booptube** | Terminal (prompts interativos) | `booptube` / `booptube.exe` |
| **booptube-gui** | Janela gráfica (tema neon futurista) | `booptube-gui` / `booptube-gui.exe` |

Ambos compartilham a mesma lógica de download, configuração e validação de URLs.

---

## O que foi implementado

### Funcionalidades comuns (CLI e GUI)

- Download de um **único vídeo** por URL (`--no-playlist`)
- Formatos **MP4** e **MP3**
- Validação de URLs do YouTube (`youtube.com`, `youtu.be`, YouTube Music, etc.)
- Escolha da **pasta de destino** (criada automaticamente se não existir)
- Arquivo salvo como `Título do vídeo.ext` (template do yt-dlp)
- **Memória da última pasta** usada (`config.json`)
- **yt-dlp** e **ffmpeg/ffprobe** embutidos e extraídos para cache na primeira execução
- Cancelamento de download em andamento

### CLI (`booptube`)

- Loop interativo: pasta → URL → formato → download → repete
- Flag `-dir` para fixar a pasta e pular o prompt
- Comandos de saída nos prompts: `q`, `sair`, `quit`, `exit`
- Progresso exibido diretamente no terminal (saída do yt-dlp)

### GUI (`booptube-gui`)

- Janela com visual **cyberpunk neon** (fundo escuro, acentos cyan e magenta)
- Campos: pasta de destino, URL, formato (MP4/MP3)
- Botão **Escolher...** para selecionar pasta no explorador de arquivos
- **Barra de progresso** com percentual parseado do yt-dlp
- **Log** com as últimas linhas do yt-dlp
- Botões **Baixar** e **Cancelar**
- Downloads consecutivos sem fechar a janela

---

## Como funciona (visão geral)

```text
Início (CLI ou GUI)
    │
    ├─ Carrega config.json (última pasta, se existir)
    ├─ Extrai yt-dlp e ffmpeg embutidos para o cache (se necessário)
    │
    └─ Interface
           ├─ CLI: loop de prompts no terminal
           └─ GUI: janela Fyne com formulário
                  │
                  └─ downloader.Download()
                         └─ subprocess yt-dlp + ffmpeg
                                └─ arquivo em %(title)s.%(ext)s
```

### Cache interno (primeira execução)

| Componente | Windows | Linux / macOS |
|------------|---------|---------------|
| yt-dlp | `%LocalAppData%\booptube\yt-dlp.exe` | `~/.cache/booptube/yt-dlp` |
| ffmpeg / ffprobe | `%LocalAppData%\booptube\ffmpeg\` | `~/.cache/booptube/ffmpeg/` |
| Configuração | `%AppData%\booptube\config.json` | `~/.config/booptube/config.json` |

A primeira execução pode demorar alguns segundos enquanto os binários embutidos são extraídos. Isso é normal.

---

## Instalação

### Opção A — Usar executável pronto (recomendado para usuários)

Se você recebeu ou baixou os arquivos compilados, basta copiá-los para uma pasta no PATH ou usar o caminho completo.

| Sistema | CLI | GUI |
|---------|-----|-----|
| Windows | `booptube.exe` | `booptube-gui.exe` |
| Linux | `booptube` | `booptube-gui` |
| macOS | `booptube` | `booptube-gui` |

**Requisitos do sistema:** nenhum extra — yt-dlp e ffmpeg já vêm embutidos.

**Tamanho aproximado:** ~185–195 MB (CLI) / ~235–245 MB (GUI). A maior parte vem do ffmpeg e yt-dlp embutidos.

**Versão:** edite [`VERSION`](../VERSION) antes do release; a CLI expõe `-version` e a GUI mostra no título.

#### Instalar no Windows (exemplo)

```powershell
# Criar pasta e copiar executáveis
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\Programs\booptube"
Copy-Item .\booptube.exe "$env:LOCALAPPDATA\Programs\booptube\"
Copy-Item .\booptube-gui.exe "$env:LOCALAPPDATA\Programs\booptube\"

# Adicionar ao PATH do usuário (opcional)
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";$env:LOCALAPPDATA\Programs\booptube",
    "User"
)
```

Depois disso, abra um novo terminal e execute `booptube` ou `booptube-gui`.

#### Instalar no Linux / macOS (exemplo)

```bash
sudo install -m 755 booptube booptube-gui /usr/local/bin/
```

---

### Opção B — Compilar a partir do código

Guias completos por sistema operacional (fetch, CLI, GUI, clean, troubleshooting):

| SO | Guia |
|----|------|
| Windows | **[build-windows.md](build-windows.md)** |
| Linux | **[build-linux.md](build-linux.md)** |
| macOS | **[build-macos.md](build-macos.md)** |

Resumo: dois binários em `./cmd/cli` (CLI) e `./cmd/gui` (GUI Fyne + CGO). Saída em `.build/`. Makefile disponível no Linux e macOS (`make build`, `make build-gui`).

---

## Como rodar

### GUI (`booptube-gui`)

#### Windows

- **Duplo clique** em `booptube-gui.exe` no Explorador de Arquivos, ou
- No terminal:

```powershell
.\booptube-gui.exe
# ou, se estiver no PATH:
booptube-gui
```

#### Linux / macOS

```bash
./booptube-gui
# ou:
booptube-gui
```

#### Passo a passo na interface

1. **Pasta de destino** — digite o caminho ou clique em **Escolher...**
2. **URL do YouTube** — cole o link do vídeo
3. **Formato** — selecione **MP4 (vídeo)** ou **MP3 (áudio)**
4. Clique em **Baixar**

Durante o download:

- A **barra de progresso** mostra o percentual (quando o yt-dlp informa)
- O **Log** exibe mensagens do yt-dlp
- **Cancelar** interrompe o download

Ao concluir, aparece **Concluído.** — você pode baixar outro vídeo sem fechar a janela.

Para sair, feche a janela. Se houver download em andamento, ele será cancelado.

### CLI (`booptube`)

#### Windows

```powershell
.\booptube.exe
.\booptube.exe -dir "C:\Downloads\booptube"
```

#### Linux / macOS

```bash
./booptube
./booptube -dir "$HOME/Downloads/booptube"
```

Sessão típica:

```text
booptube — digite q ou sair para encerrar
Pasta de destino: C:\Downloads\booptube
URL do YouTube: https://www.youtube.com/watch?v=VIDEO_ID
Formato [1=mp4, 2=mp3] (Enter=mp4):
baixando ...
concluido.
```

| Flag | Descrição |
|------|-----------|
| `-dir pasta` | Define pasta de destino e pula o prompt |
| `-h` | Ajuda |

Digite `q`, `sair`, `quit` ou `exit` nos prompts para encerrar. `Ctrl+C` cancela download ou fecha o programa.

---

## Configuração compartilhada

CLI e GUI usam o **mesmo arquivo** de configuração:

| Sistema | Caminho |
|---------|---------|
| Windows | `%AppData%\booptube\config.json` |
| Linux / macOS | `~/.config/booptube/config.json` |

Exemplo:

```json
{"download_dir":"C:\\Users\\voce\\Downloads\\booptube"}
```

A pasta é salva automaticamente após cada download bem-sucedido. Se você usar a CLI e depois abrir a GUI (ou vice-versa), a última pasta aparece pré-preenchida.

---

## Estrutura do código

```text
booptube/
├── cmd/
│   ├── cli/main.go      # Entrada booptube (CLI)
│   └── gui/main.go      # Entrada booptube-gui
├── config/              # config.json — Load/Save
├── buildinfo/           # Version, Commit, BuildDate (ldflags)
├── downloader/          # yt-dlp, ffmpeg embed, Download()
├── video/               # ParseURL, Format (mp4/mp3)
├── ui/
│   ├── terminal.go      # Interface CLI (ui.Run)
│   └── gui/             # Interface GUI Fyne (gui.Run)
├── assets/              # //go:embed yt-dlp e ffmpeg (build portable)
├── installer/           # Inno Setup, nfpm, DMG
├── scripts/             # fetch, build, stage, package
└── doc/                 # Documentação
```

### Portable vs bundled

| Variante | Tag | Embed | Uso |
|----------|-----|-------|-----|
| **Portable** | *(nenhuma)* | Sim | `make build`, zip em `.build/*-portable.*` |
| **Bundled** | `bundled` | Não — deps em `tools/` | Instaladores (`make package-*`) |

### Dois binários, um repositório

| Pasta | Comando de build | Saída |
|-------|------------------|-------|
| `cmd/cli` | `make build` ou `go build -trimpath -ldflags "-s -w ..." ./cmd/cli` | CLI |
| `cmd/gui` | `make build-gui` ou `CGO_ENABLED=1 go build ... ./cmd/gui` | GUI (Fyne) |

Versão de release: arquivo [`VERSION`](../VERSION). Builds e instaladores: **[installer.md](installer.md)**.

---

## Limitações atuais

- Apenas **um vídeo por URL** (playlists não são baixadas por completo)
- Sem seleção de qualidade/resolução (yt-dlp escolhe o melhor disponível para MP4)
- Sem download em lote não-interativo (sem passar URL por argumento)
- Nome do arquivo fixo: `%(title)s.%(ext)s`

---

## Solução de problemas

### GUI não compila: `cgo: C compiler "gcc" not found`

GUI não compila: instale MinGW (Windows) conforme [build-windows.md](build-windows.md).

### `yt-dlp embutido ausente` ou `ffmpeg embutido ausente`

Rode os scripts fetch antes de compilar:

```powershell
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
```

### `apenas URLs do YouTube sao suportadas`

Use link direto de vídeo (`watch?v=` ou `youtu.be/`).

### `pasta nao gravavel`

Escolha outra pasta ou corrija permissões de escrita.

### `download falhou`

Vídeo privado, removido, restrição regional ou problema de rede. Teste o link no navegador.

### GUI abre mas download não progride

Verifique conexão com a internet. O log na janela mostra mensagens de erro do yt-dlp.

---

## Referência rápida

| Ação | Comando |
|------|---------|
| Rodar GUI (Windows) | `booptube-gui.exe` |
| Rodar CLI (Windows) | `booptube.exe` |
| Compilar CLI | Ver [build-windows.md](build-windows.md), [build-linux.md](build-linux.md) ou [build-macos.md](build-macos.md) |
| Compilar tudo (Linux/macOS) | `make build && make build-gui` |

Para detalhes técnicos adicionais (Makefile, embed, flags), consulte [cli.md](cli.md).
Para guia de uso simplificado, consulte [usuario.md](usuario.md).
