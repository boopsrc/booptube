# Instaladores e releases

Guia para gerar **releases portable** (1 exe com embed) e **instaladores** (binários slim + pasta `tools/` compartilhada).

Visão geral: [projeto.md](projeto.md) · Build por SO: [build-windows.md](build-windows.md) · [build-linux.md](build-linux.md) · [build-macos.md](build-macos.md)

---

## Duas variantes de build

| Variante | Tag Go | Comando | Exe Windows (aprox.) | Distribuição |
|----------|--------|---------|----------------------|--------------|
| **Portable** (padrão) | *(nenhuma)* | `make build` / `scripts/build.ps1` | CLI ~214 MB, GUI ~234 MB | Zip/tar.gz com exes embed |
| **Bundled** (instalador) | `bundled` | `make build-bundled` / `scripts/build-bundled.ps1` | CLI ~4–25 MB, GUI ~35–50 MB | Setup + pasta `tools/` |

Layout instalado (bundled):

```text
{instalacao}/
├── booptube[.exe]
├── booptube-gui[.exe]
└── tools/
    ├── yt-dlp[.exe]
    ├── ffmpeg[.exe]
    └── ffprobe[.exe]
```

O app resolve `{dir-do-exe}/tools/` automaticamente. No macOS `.app`, também procura `Contents/Resources/tools/`.

---

## Artefatos em `.build/`

Binários, instaladores e zips portable ficam na mesma pasta:

| Arquivo | Tipo | Conteúdo |
|---------|------|----------|
| `booptube.exe` / `booptube-gui.exe` | Build | Executáveis compilados |
| `booptube-{v}-windows-amd64-portable.zip` | Portable | Exes com embed |
| `booptube-{v}-windows-amd64-setup.exe` | Instalador | Inno Setup, LZMA2 |
| `booptube-{v}-linux-amd64-portable.tar.gz` | Portable | Exes embed |
| `booptube-{v}-linux-amd64.deb` / `.rpm` | Instalador | nfpm |
| `booptube-{v}-macos-arm64-setup.dmg` | Instalador | `.app` + CLI em Resources |

---

## Windows

### Portable (empacotamento atual)

```powershell
.\scripts\fetch-ytdlp.ps1; .\scripts\fetch-ffmpeg.ps1
.\scripts\build.ps1
.\scripts\package-portable.ps1
```

Saída: `.build/booptube-{version}-windows-amd64-portable.zip`

### Instalador

```powershell
.\scripts\fetch-ytdlp.ps1; .\scripts\fetch-ffmpeg.ps1
.\scripts\build-bundled.ps1
.\scripts\stage.ps1 -Mode bundled
.\scripts\package.ps1
```

Pré-requisito: [Inno Setup 6](https://jrsoftware.org/isinfo.php)

```powershell
winget install JRSoftware.InnoSetup
```

Saída: `.build/booptube-{version}-windows-amd64-setup.exe`

Script Inno: [`installer/windows/booptube.iss`](../installer/windows/booptube.iss)

---

## Linux

### Portable

```bash
make fetch-deps
make build build-gui
make package-portable-linux
```

### Instalador (.deb / .rpm)

```bash
make fetch-deps
make stage          # build-bundled + staging
make package-linux  # requer nfpm no PATH
```

Instalar nfpm: [github.com/goreleaser/nfpm](https://github.com/goreleaser/nfpm)

```bash
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

Instalar o .deb:

```bash
sudo dpkg -i .build/booptube_*_amd64.deb
booptube-gui
```

**Sem nfpm:** `package-linux` gera tarball bundled; extraia e rode `sudo installer/linux/install.sh`.

---

## macOS

### Portable

```bash
make fetch-deps
make build build-gui
make package-portable-macos
```

### Instalador (DMG)

```bash
make fetch-deps
make stage
make package-macos
```

Arraste `booptube-gui.app` para Applications. CLI fica em `booptube-gui.app/Contents/Resources/booptube`.

Se o macOS bloquear:

```bash
xattr -cr /Applications/booptube-gui.app
```

---

## CI / Release automático (GitHub Actions)

O workflow [`.github/workflows/release.yml`](../.github/workflows/release.yml) publica um **GitHub Release** quando há push em `main` com mensagem de commit contendo **`Bump version`**.

### Fluxo do mantenedor

```bash
echo "0.2.0" > VERSION
git add VERSION
git commit -m "Bump version to 0.2.0"
git push origin main
```

1. O job `gate` lê a versão em [`VERSION`](../VERSION)
2. Três jobs compilam **em paralelo**: Windows, Linux, macOS
3. O job `release` agrupa os artefatos e cria a tag `v{VERSION}` (ex.: `v0.2.0`)

### Artefatos no GitHub Release

| Plataforma | Portable | Instalador |
|------------|----------|------------|
| Windows | `booptube-{v}-windows-amd64-portable.zip` | `booptube-{v}-windows-amd64-setup.exe` |
| Linux | `booptube-{v}-linux-amd64-portable.tar.gz` | `booptube_{v}_amd64.deb`, `.rpm` |
| macOS | `booptube-{v}-macos-arm64-portable.tar.gz` | `booptube-{v}-macos-arm64-setup.dmg` |

Cada bump exige **versão nova** em `VERSION` — se a tag `v{VERSION}` já existir, o workflow falha.

---

## Makefile (referência)

| Target | Descrição |
|--------|-----------|
| `build` / `build-gui` | Portable (embed) — **padrão, inalterado** |
| `build-bundled` / `build-gui-bundled` | Slim para instalador |
| `stage` | Bundled + `tools/` em `installer/staging/` |
| `package-portable` | Zip/tar portable do SO atual |
| `package-portable-win` / `-linux` / `-macos` | Portable por SO |
| `package` | Instalador do SO atual |
| `package-win` / `-linux` / `-macos` | Instalador por SO |
| `clean` | Remove `.build/` e `installer/staging/` |

---

## Solução de problemas

| Problema | Solução |
|----------|---------|
| `yt-dlp nao encontrado` (bundled) | Reinstale; verifique `{app}/tools/` |
| `yt-dlp embutido ausente` (portable) | Rode fetch-deps antes do build |
| `ISCC not found` | Instale Inno Setup |
| `nfpm not found` | Instale nfpm ou use tarball + `install.sh` |
| Release CI falhou (tag existe) | Use versão nova em `VERSION` antes do commit Bump version |
| GUI abre terminal (Windows) | Recompile GUI com `-H=windowsgui` (já no `build.ps1`) |

---

## Ver também

- [projeto.md](projeto.md) — arquitetura portable vs bundled
- [build-windows.md](build-windows.md) — compilar no Windows
