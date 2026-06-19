# booptube

CLI para baixar videos do YouTube em mp4 ou mp3.

Documentacao completa: [doc/cli.md](doc/cli.md)

## Build

```powershell
.\fetch-ytdlp.ps1
go build -o .build/booptube.exe .
```

Linux/macOS: `make fetch-ytdlp && make build` → `./.build/booptube`

## Uso

```text
.\.build\booptube.exe
.\.build\booptube.exe -dir C:\Downloads
```

No loop interativo: pasta de destino, URL do YouTube, formato (1=mp4, 2=mp3). Digite `q` ou `sair` para encerrar.

A ultima pasta e salva em `%AppData%\booptube\config.json` (Windows) ou `~/.config/booptube/config.json`.

## Dependencias

- **ffmpeg** no PATH — necessario para merge mp4 e conversao mp3 (`winget install Gyan.FFmpeg`)
