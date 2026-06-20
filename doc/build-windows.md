# Compilar no Windows

Guia passo a passo para build da **CLI** (`booptube.exe`) e **GUI** (`booptube-gui.exe`) no Windows.

Fluxo recomendado: **PowerShell** + scripts `.ps1`. O comando `make` não vem instalado por padrão no Windows.

Outros guias: [build-linux.md](build-linux.md) · [build-macos.md](build-macos.md)

---

## Pré-requisitos

| Requisito | Versão | CLI | GUI |
|-----------|--------|-----|-----|
| [Go](https://go.dev/dl/) | 1.22+ | Sim | Sim |
| PowerShell | 5+ | Sim | Sim |
| GCC (MinGW-w64) | — | Não | Sim (CGO / Fyne) |

**yt-dlp** e **ffmpeg** não precisam estar instalados no sistema — são embutidos no executável após o fetch.

### Instalar GCC para a GUI

Escolha uma opção:

**MSYS2 (recomendado):**

1. Instale [MSYS2](https://www.msys2.org/)
2. No terminal **MSYS2 MinGW x64**:

```bash
pacman -S mingw-w64-x86_64-gcc
```

3. Adicione ao PATH do Windows: `C:\msys64\mingw64\bin`
4. Abra um **novo** PowerShell e confirme:

```powershell
gcc --version
```

**WinLibs (via winget):**

```powershell
winget install -e --id BrechtSanders.WinLibs.POSIX.UCRT
```

Reabra o terminal após a instalação.

**TDM-GCC:** [jmeubank.github.io/tdm-gcc](https://jmeubank.github.io/tdm-gcc/)

---

## 1. Baixar yt-dlp e ffmpeg

Obrigatório na **primeira compilação** ou após atualizar versões embutidas.

```powershell
cd C:\caminho\para\booptube
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
```

Os binários ficam em `assets/ytdlp/` e `assets/ffmpeg/` (não vão para o git).

---

## 2. Compilar (portable — padrão)

Após o fetch, use o script de build com flags otimizadas (`-s -w`, versão injetada, GUI sem console). **Sem tag `bundled`** — yt-dlp e ffmpeg ficam embutidos no exe (empacotamento atual).

```powershell
.\scripts\build.ps1
```

Compilar só a CLI ou só a GUI:

```powershell
.\scripts\build.ps1 -Target cli
.\scripts\build.ps1 -Target gui
```

Saída:
- `.build/booptube.exe` (~185–195 MB)
- `.build/booptube-gui.exe` (~235–245 MB)

> A maior parte do tamanho vem do **ffmpeg e yt-dlp embutidos** (~90–130 MB). As flags `-s -w` reduzem só a parte Go/Fyne (~5–15 MB).

### Compilar manualmente

```powershell
New-Item -ItemType Directory -Force -Path .build
go build -trimpath -ldflags "-s -w" -o .build/booptube.exe ./cmd/cli
$env:CGO_ENABLED = "1"
go build -trimpath -ldflags "-s -w -H=windowsgui" -o .build/booptube-gui.exe ./cmd/gui
```

A flag `-H=windowsgui` evita que uma janela de terminal abra junto com a GUI no duplo clique.

A primeira compilação da GUI pode demorar vários minutos (Fyne + CGO).

---

## 3. Versão

A versão de release fica em [`VERSION`](../VERSION) na raiz. O build injeta versão, commit e data via `-ldflags`.

```powershell
.\.build\booptube.exe -version
# booptube 0.1.0 (abc1234, 2026-06-19T...)
```

A GUI exibe a versão no título da janela e no subtítulo da interface.

---

## 4. Limpar build

```powershell
Remove-Item -Recurse -Force .build -ErrorAction SilentlyContinue
```

---

## 5. Rodar

```powershell
.\.build\booptube.exe
.\.build\booptube.exe -dir "C:\Downloads\booptube"
.\.build\booptube.exe -version
.\.build\booptube-gui.exe
```

A GUI abre **somente a janela gráfica** (sem terminal) no duplo clique no Explorador de Arquivos.

### Adicionar ao PATH (opcional)

```powershell
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\Programs\booptube"
Copy-Item .\.build\booptube.exe "$env:LOCALAPPDATA\Programs\booptube\"
Copy-Item .\.build\booptube-gui.exe "$env:LOCALAPPDATA\Programs\booptube\"
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";$env:LOCALAPPDATA\Programs\booptube",
    "User"
)
```

---

## Problemas comuns

| Erro | Solução |
|------|---------|
| `yt-dlp embutido ausente` | Rode `fetch-ytdlp.ps1` e `fetch-ffmpeg.ps1`, depois recompile |
| `cgo: C compiler "gcc" not found` | Instale MinGW, adicione ao PATH, abra novo PowerShell |
| `make` não reconhecido | Normal no PowerShell — use os comandos deste guia |
| Build da GUI muito lento | Normal na 1ª vez; compilações seguintes são mais rápidas |

---

## Usar `make` no Windows (opcional)

O [`Makefile`](../Makefile) chama scripts **`.sh`** (bash). No PowerShell puro, `make build` **não funciona** mesmo com `make` instalado.

| Opção | Instalação | Terminal para `make` |
|-------|------------|----------------------|
| MSYS2 | `pacman -S make` | MSYS2 MinGW x64 |
| Chocolatey | `choco install make` | Git Bash ou MSYS2 |
| Git for Windows | incluído no Git Bash | Git Bash |

No **Git Bash** ou **MSYS2**, na pasta do projeto:

```bash
make build
make build-gui
```

No dia a dia no Windows, prefira **PowerShell + scripts `.ps1`** (seções 1–3 acima).

---

## Referência Makefile

| Target | Equivalente PowerShell |
|--------|------------------------|
| `make fetch-deps` | `.\scripts\fetch-ytdlp.ps1; .\scripts\fetch-ffmpeg.ps1` |
| `make build` | fetch + `.\scripts\build.ps1 -Target cli` |
| `make build-gui` | fetch + `.\scripts\build.ps1 -Target gui` |
| build completo | `.\scripts\build.ps1` |
| `make clean` | `Remove-Item -Recurse -Force .build` |

---

## Instalador e release portable

| Tipo | Comandos | Saída |
|------|----------|-------|
| **Portable zip** | `.\scripts\build.ps1` → `.\scripts\package-portable.ps1` | `.build/*-portable.zip` |
| **Instalador slim** | `.\scripts\build-bundled.ps1` → `.\scripts\stage.ps1` → `.\scripts\package.ps1` | `.build/*-setup.exe` |

Guia completo: **[installer.md](installer.md)** (Inno Setup, tamanhos slim vs portable).

---

## Ver também

- [installer.md](installer.md) — instaladores e releases

- [gui.md](gui.md) — usar a interface gráfica
- [cli.md](cli.md) — referência técnica e Makefile
- [projeto.md](projeto.md) — visão geral do projeto
