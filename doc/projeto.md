# booptube — documentação do projeto

Visão geral de **tudo que foi implementado** até o momento: o que o booptube faz, como está organizado, como instalar e como rodar a **CLI** e a **GUI**.

| Guia relacionado | Conteúdo |
|------------------|----------|
| [usuario.md](usuario.md) | Uso do dia a dia (CLI e GUI) |
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

**Tamanho aproximado:** ~200 MB por executável (dependências embutidas por plataforma).

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

#### Pré-requisitos gerais

| Requisito | Versão | Necessário para |
|-----------|--------|-----------------|
| Go | 1.22+ | Compilar CLI e GUI |
| PowerShell | 5+ | Scripts fetch no Windows |
| bash, curl, unzip | — | Scripts fetch no Linux/macOS |

#### Compilar a CLI

**Windows:**

```powershell
cd booptube
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
go build -o .build/booptube.exe .
```

**Linux / macOS:**

```bash
cd booptube
make build
# ou manualmente:
chmod +x scripts/*.sh
./scripts/fetch-ytdlp.sh
./scripts/fetch-ffmpeg.sh
go build -o .build/booptube .
```

Resultado: `.build/booptube` ou `.build/booptube.exe`

> Os binários em `assets/ytdlp/` e `assets/ffmpeg/` não estão no git. Rode os scripts fetch antes de compilar.

#### Compilar a GUI

A GUI usa [Fyne](https://fyne.io/) e exige **CGO habilitado** (`CGO_ENABLED=1`). No Windows, é necessário um compilador C (GCC).

##### Pré-requisitos extras da GUI

| Sistema | Requisito | Como instalar |
|---------|-----------|---------------|
| **Windows** | GCC (MinGW-w64) | [MSYS2](https://www.msys2.org/): `pacman -S mingw-w64-x86_64-gcc` — adicione `C:\msys64\mingw64\bin` ao PATH |
| **Windows** | Alternativa | [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) |
| **Linux** | pkg-config, libgl, gcc | Debian/Ubuntu: `sudo apt install gcc libgl1-mesa-dev xorg-dev pkg-config` |
| **macOS** | Xcode CLT | `xcode-select --install` |

##### Verificar se o GCC está disponível (Windows)

```powershell
gcc --version
# deve mostrar a versão do MinGW, ex.: x86_64-w64-mingw32-gcc
```

Se aparecer `gcc not found`, instale o MinGW e reinicie o terminal.

##### Comandos de build da GUI

**Windows:**

```powershell
cd booptube
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
$env:CGO_ENABLED = "1"
go build -tags gui -o .build/booptube-gui.exe .
```

**Linux / macOS:**

```bash
cd booptube
make build-gui
# ou:
CGO_ENABLED=1 go build -tags gui -o .build/booptube-gui .
```

Resultado: `.build/booptube-gui` ou `.build/booptube-gui.exe`

##### Makefile — targets disponíveis

| Comando | Saída |
|---------|-------|
| `make build` | `.build/booptube` (CLI) |
| `make build-gui` | `.build/booptube-gui` (GUI, requer CGO) |
| `make fetch-deps` | Baixa yt-dlp e ffmpeg para embed |
| `make clean` | Remove `.build/` |

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
├── main.go              # Entrada CLI (build tag !gui)
├── main_gui.go          # Entrada GUI (build tag gui)
├── config/              # config.json — Load/Save
├── downloader/          # yt-dlp, ffmpeg embed, Download()
├── video/               # ParseURL, Format (mp4/mp3)
├── ui/
│   ├── terminal.go      # Interface CLI
│   ├── gui.go           # Interface GUI (Fyne)
│   └── gui_theme.go     # Tema neon
├── assets/              # //go:embed yt-dlp e ffmpeg por OS
├── scripts/             # fetch-ytdlp, fetch-ffmpeg
└── doc/                 # Documentação
```

### Dois binários, um repositório

| Arquivo | Build tag | Comando de build |
|---------|-----------|------------------|
| `main.go` | `!gui` | `go build -o booptube .` |
| `main_gui.go` | `gui` | `go build -tags gui -o booptube-gui .` |

A tag `gui` inclui os arquivos `ui/gui.go` e `ui/gui_theme.go`, que dependem do Fyne.

---

## Limitações atuais

- Apenas **um vídeo por URL** (playlists não são baixadas por completo)
- Sem seleção de qualidade/resolução (yt-dlp escolhe o melhor disponível para MP4)
- Sem download em lote não-interativo (sem passar URL por argumento)
- Nome do arquivo fixo: `%(title)s.%(ext)s`

---

## Solução de problemas

### GUI não compila: `cgo: C compiler "gcc" not found`

Instale MinGW-w64 (MSYS2) no Windows e garanta que `gcc` está no PATH. Veja a seção [Compilar a GUI](#compilar-a-gui).

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
| Compilar CLI | `go build -o .build/booptube.exe .` |
| Compilar GUI (Windows) | `$env:CGO_ENABLED="1"; go build -tags gui -o .build/booptube-gui.exe .` |
| Compilar tudo (Linux/macOS) | `make build && make build-gui` |

Para detalhes técnicos adicionais (Makefile, embed, flags), consulte [cli.md](cli.md).
Para guia de uso simplificado, consulte [usuario.md](usuario.md).
