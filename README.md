# booptube

CLI para baixar vídeos do YouTube em **mp4** ou **mp3**. Embute **yt-dlp** e **ffmpeg** — nada extra para instalar além do executável.

## Documentação

| Guia | Para quem |
|------|-----------|
| **[doc/usuario.md](doc/usuario.md)** | Usuário — como usar a CLI no dia a dia |
| **[doc/cli.md](doc/cli.md)** | Desenvolvedor — build, config técnica, Makefile |
| **[doc/README.md](doc/README.md)** | Índice da documentação |

## Uso rápido

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

```bash
./scripts/fetch-ytdlp.sh
./scripts/fetch-ffmpeg.sh
go build -o .build/booptube .
```

Windows:

```powershell
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
go build -o .build/booptube.exe .
```

Linux/macOS: `make build` → `./.build/booptube`

Detalhes em [doc/cli.md](doc/cli.md).
