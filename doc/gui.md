# booptube-gui — instalação e uso

Guia focado na **interface gráfica** do booptube: como instalar, compilar e usar o `booptube-gui`.

Visão geral do projeto: [projeto.md](projeto.md) · Uso geral: [usuario.md](usuario.md) · Build técnico: [cli.md](cli.md)

---

## O que é o booptube-gui

Aplicativo com **janela gráfica** para baixar vídeos do YouTube em MP4 ou MP3. Visual **neon futurista** (fundo escuro, acentos cyan e magenta).

Faz tudo que a CLI faz, sem terminal:

- Escolher pasta de destino (digitando ou pelo explorador)
- Colar URL do YouTube
- Escolher MP4 ou MP3
- Ver progresso e log do download
- Cancelar ou baixar vários vídeos seguidos

**Não precisa instalar yt-dlp nem ffmpeg** — vêm embutidos no executável.

---

## Instalação

### 1. Executável pronto (mais simples)

Se você já tem `booptube-gui.exe` (Windows) ou `booptube-gui` (Linux/macOS):

1. Copie o arquivo para uma pasta de sua preferência
2. Execute com duplo clique (Windows) ou `./booptube-gui` (Linux/macOS)

**Opcional — atalho no Windows:**

- Clique direito em `booptube-gui.exe` → **Fixar na barra de tarefas** ou **Criar atalho**

**Opcional — adicionar ao PATH (Windows):**

```powershell
Copy-Item .\booptube-gui.exe "$env:LOCALAPPDATA\Programs\booptube\"
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";$env:LOCALAPPDATA\Programs\booptube",
    "User"
)
```

Reabra o terminal. Agora `booptube-gui` funciona de qualquer pasta.

---

### 2. Compilar você mesmo

Necessário se você clonou o repositório e ainda não tem o executável.

#### Passo 1 — Instalar Go

Baixe em [go.dev/dl](https://go.dev/dl/) (versão **1.22+**).

#### Passo 2 — Instalar compilador C (obrigatório para a GUI)

A GUI usa Fyne, que depende de **CGO**.

**Windows — MSYS2 + MinGW (recomendado):**

1. Instale [MSYS2](https://www.msys2.org/)
2. Abra **MSYS2 MinGW x64** e execute:

```bash
pacman -S mingw-w64-x86_64-gcc
```

3. Adicione ao PATH do Windows:

```text
C:\msys64\mingw64\bin
```

4. Abra um **novo** PowerShell e confirme:

```powershell
gcc --version
```

**Linux (Debian/Ubuntu):**

```bash
sudo apt update
sudo apt install gcc libgl1-mesa-dev xorg-dev pkg-config
```

**macOS:**

```bash
xcode-select --install
```

#### Passo 3 — Baixar dependências embutidas (yt-dlp + ffmpeg)

**Windows:**

```powershell
cd booptube
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
```

**Linux / macOS:**

```bash
cd booptube
make fetch-deps
```

#### Passo 4 — Compilar a GUI

**Windows:**

```powershell
$env:CGO_ENABLED = "1"
go build -o .build/booptube-gui.exe ./cmd/gui
```

**Linux / macOS:**

```bash
make build-gui
```

O executável ficará em:

| Sistema | Caminho |
|---------|---------|
| Windows | `.build/booptube-gui.exe` |
| Linux / macOS | `.build/booptube-gui` |

#### Erro comum na compilação

```text
cgo: C compiler "gcc" not found
```

**Solução:** instale MinGW (Windows) ou gcc (Linux) conforme o Passo 2 e reinicie o terminal.

---

## Como rodar

### Windows

```powershell
# A partir da pasta do projeto (após compilar)
.\.build\booptube-gui.exe

# Ou, se copiou para outro lugar:
C:\caminho\para\booptube-gui.exe
```

Também funciona **duplo clique** no `.exe`.

### Linux / macOS

```bash
./.build/booptube-gui
```

---

## Usando a interface

Ao abrir, você verá:

```text
┌─────────────────────────────────────┐
│  booptube          (título neon)    │
├─────────────────────────────────────┤
│  Pasta de destino                   │
│  [________________________] Escolher│
│  URL do YouTube                     │
│  [________________________]         │
│  Formato                            │
│  ( ) MP4 (vídeo)  ( ) MP3 (áudio)   │
│  [ Baixar ]  [ Cancelar ]           │
│  Status: Pronto.                    │
│  [======== progresso ========]      │
│  Log                                │
│  [                              ]   │
└─────────────────────────────────────┘
```

### Fluxo

1. **Pasta de destino**
   - Digite o caminho (ex.: `C:\Users\voce\Downloads`)
   - Ou clique **Escolher...** e selecione no explorador
   - A última pasta usada é carregada automaticamente

2. **URL do YouTube**
   - Cole o link do vídeo, ex.:
     - `https://www.youtube.com/watch?v=VIDEO_ID`
     - `https://youtu.be/VIDEO_ID`

3. **Formato**
   - **MP4 (vídeo)** — melhor qualidade de vídeo + áudio
   - **MP3 (áudio)** — só áudio convertido

4. Clique **Baixar**

### Durante o download

- Campos ficam desabilitados
- **Cancelar** aparece — clique para interromper
- Barra de progresso avança quando o yt-dlp informa percentual
- **Log** mostra mensagens em tempo real

### Após concluir

- Mensagem **Concluído.** no status
- URL é limpa; pasta permanece
- Você pode colar outro link e baixar de novo

### Onde fica o arquivo

Na pasta escolhida, com o nome do título do vídeo:

```text
C:\Downloads\Me at the zoo.mp4
```

---

## Primeira execução

Na **primeira vez** (ou após atualizar o booptube), o programa extrai yt-dlp e ffmpeg para o cache do usuário. Pode levar alguns segundos antes da janela responder normalmente.

| Sistema | Cache |
|---------|-------|
| Windows | `%LocalAppData%\booptube\` |
| Linux / macOS | `~/.cache/booptube/` |

---

## Configuração

A GUI salva a última pasta em:

| Sistema | Arquivo |
|---------|---------|
| Windows | `%AppData%\booptube\config.json` |
| Linux / macOS | `~/.config/booptube/config.json` |

CLI e GUI compartilham o mesmo arquivo — a pasta escolhida na GUI aparece na CLI e vice-versa.

---

## Problemas comuns

| Problema | O que fazer |
|----------|-------------|
| Janela não abre | Execute pelo terminal para ver mensagens de erro |
| `gcc not found` ao compilar | Instale MinGW e adicione ao PATH |
| `Informe a pasta de destino` | Preencha o campo ou use Escolher... |
| `apenas URLs do YouTube` | Use link de vídeo do YouTube, não de outro site |
| `pasta nao gravavel` | Escolha pasta com permissão de escrita (ex.: Downloads) |
| Progresso parado em 0% | Normal em algumas fases; confira o Log |
| Download falhou | Vídeo privado/removido ou sem internet — teste no navegador |

---

## Diferença entre CLI e GUI

| | CLI (`booptube`) | GUI (`booptube-gui`) |
|--|------------------|----------------------|
| Interface | Terminal | Janela gráfica |
| Compilar | Só Go | Go + GCC (CGO) |
| Flag `-dir` | Sim | Não (usa campo na tela) |
| Progresso | Texto no terminal | Barra + log |
| Cancelar | Ctrl+C | Botão Cancelar |

Ambos produzem o **mesmo arquivo** na pasta escolhida.
