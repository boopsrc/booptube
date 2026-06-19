# booptube — documentação técnica da CLI

Referência para **desenvolvedores** e quem compila o projeto. Para uso do dia a dia, veja **[usuario.md](usuario.md)**.

CLI interativa para baixar vídeos do YouTube em **mp4** ou **mp3**. O binário embute **yt-dlp** e **ffmpeg** (essentials); na primeira execução são extraídos para o cache do usuário.

## Índice

1. [Como funciona](#como-funciona)
2. [Instalação e build](#instalação-e-build)
3. [Configuração](#configuração)
4. [Uso da CLI](#uso-da-cli)
5. [Comandos e flags](#comandos-e-flags)
6. [Prompts interativos](#prompts-interativos)
7. [Comandos do Makefile](#comandos-do-makefile)
8. [Arquivos gerados](#arquivos-gerados)
9. [Dependências embutidas](#dependências-embutidas)
10. [Solução de problemas](#solução-de-problemas)

---

## Como funciona

```text
booptube inicia
    │
    ├─ Carrega config.json (última pasta, se existir)
    ├─ Extrai yt-dlp e ffmpeg embutidos para o cache (se necessário)
    │
    └─ Loop interativo
           ├─ Pergunta pasta de destino (ou usa -dir / config salva)
           ├─ Pergunta URL do YouTube
           ├─ Pergunta formato (mp4 ou mp3)
           ├─ Executa yt-dlp e mostra progresso no terminal
           └─ Salva a pasta em config.json e repete
```

Fluxo resumido:

1. Você informa **onde** salvar, **qual** vídeo e **em qual formato**.
2. O booptube chama o **yt-dlp** (incluído no executável) com **ffmpeg** embutido via `--ffmpeg-location`.
3. O arquivo é gravado na pasta escolhida com o nome `Título do vídeo.ext`.
4. Após cada download bem-sucedido, a **última pasta** é memorizada para a próxima sessão.

**Cancelamento:** `Ctrl+C` interrompe o download em andamento. Nos prompts, digite `q`, `sair`, `quit` ou `exit` para encerrar o programa.

**Playlists:** apenas o vídeo da URL informada é baixado (`--no-playlist`).

---

## Instalação e build

### Pré-requisitos

| Requisito | Versão mínima | Observação |
|-----------|---------------|------------|
| Go | 1.22+ | só para compilar a partir do código |
| bash / curl / unzip | Linux e macOS | para `scripts/*.sh` e `make build` |
| PowerShell | 5+ | para `scripts/*.ps1` no Windows |

**ffmpeg e yt-dlp não precisam ser instalados no sistema** — vêm embutidos no executável compilado.

### Compilar a partir do código

#### Windows

```powershell
cd booptube
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
go build -o .build/booptube.exe .
```

O executável fica em `.build/booptube.exe` (~200 MB).

#### Linux / macOS

```bash
cd booptube
chmod +x scripts/*.sh
./scripts/fetch-ytdlp.sh
./scripts/fetch-ffmpeg.sh
go build -o .build/booptube .
```

Ou, em um comando: `make build`.

O executável fica em `.build/booptube`.

> **Nota:** os binários em `assets/ytdlp/` e `assets/ffmpeg/` não vão para o git. Rode os scripts fetch ou `make build` antes de compilar.

### Adicionar ao PATH (opcional)

Copie ou crie um atalho do executável em uma pasta que já esteja no PATH, ou invoque pelo caminho completo:

```powershell
# Windows — exemplo
Copy-Item .\.build\booptube.exe "$env:LOCALAPPDATA\Programs\booptube\"
```

---

## Configuração

### Arquivo de configuração

Após o primeiro download concluído, o booptube grava:

| Sistema | Caminho |
|---------|---------|
| Windows | `%AppData%\booptube\config.json` |
| Linux | `~/.config/booptube/config.json` |
| macOS | `~/.config/booptube/config.json` |

Exemplo de conteúdo:

```json
{"download_dir":"C:\\Users\\voce\\Downloads\\booptube"}
```

Campos:

| Campo | Descrição |
|-------|-----------|
| `download_dir` | Última pasta de destino usada com sucesso |

Não é necessário editar manualmente: a CLI atualiza o arquivo ao fim de cada download. Você pode editar o JSON para fixar uma pasta padrão antes de abrir o programa.

### Cache do yt-dlp

Na primeira execução, o yt-dlp embutido é copiado para:

| Sistema | Caminho |
|---------|---------|
| Windows | `%LocalAppData%\booptube\yt-dlp.exe` |
| Linux / macOS | `~/.cache/booptube/yt-dlp` |

Se você recompilar o booptube com uma versão nova do yt-dlp, o cache é atualizado automaticamente (checksum diferente).

### Cache do ffmpeg

Na primeira execução, **ffmpeg** e **ffprobe** embutidos são copiados para:

| Sistema | Caminho |
|---------|---------|
| Windows | `%LocalAppData%\booptube\ffmpeg\ffmpeg.exe` e `ffprobe.exe` |
| Linux / macOS | `~/.cache/booptube/ffmpeg/ffmpeg` e `ffprobe` |

O yt-dlp recebe `--ffmpeg-location` apontando para essa pasta. Recompilar com ffmpeg novo atualiza o cache automaticamente.

### Variáveis de ambiente (build)

| Variável | Uso |
|----------|-----|
| `YTDLP_VERSION` | Versão do yt-dlp em `scripts/fetch-ytdlp.{sh,ps1}` / `make fetch-ytdlp` (padrão: `2026.06.09`) |
| `FFMPEG_VERSION` | Versão Gyan essentials no Windows (padrão: `8.1.1`) |

Exemplo:

```bash
export YTDLP_VERSION=2026.06.09
export FFMPEG_VERSION=8.1.1
./scripts/fetch-ytdlp.sh
./scripts/fetch-ffmpeg.sh
```

---

## Uso da CLI

Guia passo a passo para usuários finais: **[usuario.md](usuario.md)**.

### Modo interativo (padrão)

```powershell
# Windows
.\.build\booptube.exe
```

```bash
# Linux / macOS
./.build/booptube
```

Sessão típica:

```text
booptube — digite q ou sair para encerrar
Pasta de destino: C:\Downloads\booptube
URL do YouTube: https://www.youtube.com/watch?v=VIDEO_ID
Formato [1=mp4, 2=mp3] (Enter=mp4): 1
baixando https://www.youtube.com/watch?v=VIDEO_ID como mp4...
[progresso do yt-dlp]
concluido.
Pasta de destino (Enter=C:\Downloads\booptube):
```

Na segunda rodada do loop, pressionar **Enter** na pasta reutiliza o caminho anterior.

### Com pasta fixa (flag `-dir`)

Pula o prompt da pasta em todas as iterações:

```powershell
.\.build\booptube.exe -dir "C:\Downloads\booptube"
```

```bash
./.build/booptube -dir "$HOME/Downloads/booptube"
```

A pasta é criada automaticamente se não existir. Deve ser gravável.

---

## Comandos e flags

### Executável `booptube`

```text
booptube [flags]
```

| Flag | Tipo | Padrão | Descrição |
|------|------|--------|-----------|
| `-dir` | string | *(vazio)* | Pasta de destino. Quando informada, o prompt da pasta é omitido. |

**Ajuda:**

```bash
go run . -h
# ou, após compilar:
./.build/booptube -h
```

Saída esperada:

```text
  -dir string
        pasta de destino (pula prompt da pasta)
```

Não há outros flags na versão atual.

### Comandos de saída (dentro dos prompts)

Digite em qualquer prompt de texto:

| Comando | Ação |
|---------|------|
| `q` | Encerra o programa |
| `quit` | Encerra o programa |
| `sair` | Encerra o programa |
| `exit` | Encerra o programa |

### Atalhos de teclado

| Tecla | Ação |
|-------|------|
| `Ctrl+C` | Cancela download ou encerra o processo |
| `Enter` | Confirma valor padrão (pasta anterior ou formato mp4) |

---

## Prompts interativos

### 1. Pasta de destino

```text
Pasta de destino: 
Pasta de destino (Enter=C:\caminho\anterior):
```

- Aceita caminho absoluto ou relativo.
- Cria a pasta se não existir.
- Falha se a pasta não for gravável.
- Com `-dir`, este prompt não aparece.

### 2. URL do YouTube

```text
URL do YouTube:
```

URLs aceitas (hosts):

- `youtube.com`, `www.youtube.com`, `m.youtube.com`, `music.youtube.com`
- `youtu.be`
- Subdomínios `*.youtube.com`

O esquema `https://` é adicionado automaticamente se omitido.

Exemplos válidos:

```text
https://www.youtube.com/watch?v=dQw4w9WgXcQ
youtu.be/dQw4w9WgXcQ
www.youtube.com/watch?v=dQw4w9WgXcQ
```

### 3. Formato

```text
Formato [1=mp4, 2=mp3] (Enter=mp4):
```

| Entrada | Resultado |
|---------|-----------|
| *(Enter vazio)* | mp4 |
| `1`, `mp4`, `video` | mp4 |
| `2`, `mp3`, `audio` | mp3 |

**mp4:** melhor vídeo + áudio, merge em mp4 (ffmpeg embutido).

**mp3:** extrai áudio e converte para mp3 (ffmpeg embutido).

---

## Comandos do Makefile

Disponíveis em Linux e macOS (requer `bash`, `curl`, `unzip` e `tar`):

| Comando | Descrição |
|---------|-----------|
| `make fetch-ytdlp` | Executa `scripts/fetch-ytdlp.sh` |
| `make fetch-ffmpeg` | Executa `scripts/fetch-ffmpeg.sh` |
| `make fetch-deps` | Roda os dois scripts acima |
| `make build` | Roda `fetch-deps` e compila para `.build/booptube` |
| `make clean` | Remove a pasta `.build/` |

Versões pinadas: `YTDLP_VERSION ?= 2026.06.09`, `FFMPEG_VERSION ?= 8.1.1`.

Override:

```bash
make build YTDLP_VERSION=2026.06.09 FFMPEG_VERSION=8.1.1
```

### Scripts shell (Linux / macOS)

| Script | Descrição |
|--------|-----------|
| `./scripts/fetch-ytdlp.sh` | Equivalente ao `make fetch-ytdlp` |
| `./scripts/fetch-ffmpeg.sh` | Equivalente ao `make fetch-ffmpeg` |

### Scripts PowerShell (Windows)

| Script | Descrição |
|--------|-----------|
| `.\scripts\fetch-ytdlp.ps1` | Equivalente ao `make fetch-ytdlp` |
| `.\scripts\fetch-ffmpeg.ps1` | Equivalente ao `make fetch-ffmpeg` |

---

## Arquivos gerados

### Na pasta de destino

| Formato | Nome típico |
|---------|-------------|
| mp4 | `Título do vídeo.mp4` |
| mp3 | `Título do vídeo.mp3` |

O nome vem do título do vídeo no YouTube (`%(title)s` do yt-dlp).

### Pastas do projeto (desenvolvimento)

| Caminho | Conteúdo |
|---------|------------|
| `.build/` | Executável compilado (ignorado pelo git) |
| `assets/ytdlp/` | Binários yt-dlp antes do embed (ignorado pelo git) |
| `assets/ffmpeg/` | Binários ffmpeg/ffprobe antes do embed (ignorado pelo git) |

---

## Dependências embutidas

| Componente | Incluído no booptube? | Necessário? |
|------------|----------------------|-------------|
| yt-dlp | Sim (embutido) | Automático na extração |
| ffmpeg + ffprobe | Sim (embutido, essentials) | Automático na extração |
| Go | Não | Só para compilar |

Tamanho aproximado do executável: **~200 MB** (yt-dlp + ffmpeg/ffprobe essentials por plataforma).

---

## Solução de problemas

### `yt-dlp embutido ausente: execute fetch-ytdlp...`

O build foi feito sem baixar os assets. Rode `scripts/fetch-ytdlp.sh`, `scripts/fetch-ytdlp.ps1` ou `make fetch-ytdlp` e recompile.

### `apenas URLs do YouTube sao suportadas`

A URL não é de um host YouTube reconhecido. Use link direto de vídeo (`watch?v=` ou `youtu.be/`).

### `pasta nao gravavel`

Escolha outra pasta ou corrija permissões. O booptube testa escrita antes de baixar.

### `ffmpeg embutido ausente: execute fetch-ffmpeg...`

O build foi feito sem baixar ffmpeg. Rode `scripts/fetch-ffmpeg.sh`, `scripts/fetch-ffmpeg.ps1` ou `make fetch-ffmpeg` e recompile.

### `Postprocessing: ffprobe and ffmpeg not found`

Cache corrompido ou extração incompleta. Apague `%LocalAppData%\booptube\ffmpeg\` (ou `~/.cache/booptube/ffmpeg/`) e execute o booptube novamente.

### Download lento ou falha de rede

- Verifique conexão.
- URLs muito longas ou vídeos restritos podem falhar no yt-dlp.
- Recompile com yt-dlp mais recente: `make fetch-ytdlp` + `make build`.

### Atualizar ffmpeg embutido

1. Ajuste `FFMPEG_VERSION` em `Makefile` / `scripts/fetch-ffmpeg.{sh,ps1}` se quiser pinar outra release.
2. Rode o fetch e recompile.
3. Na próxima execução, o cache local do ffmpeg será substituído.

### Atualizar yt-dlp embutido

1. Ajuste `YTDLP_VERSION` em `Makefile` / `scripts/fetch-ytdlp.{sh,ps1}` se quiser pinar outra release.
2. Rode o fetch e recompile.
3. Na próxima execução, o cache local do yt-dlp será substituído.

---

## Referência rápida

**Usuário:** [usuario.md](usuario.md)

**Desenvolvedor:**

```bash
# Compilar (Linux/macOS)
chmod +x scripts/*.sh
./scripts/fetch-ytdlp.sh
./scripts/fetch-ffmpeg.sh
go build -o .build/booptube .

# Compilar (Windows)
.\scripts\fetch-ytdlp.ps1
.\scripts\fetch-ffmpeg.ps1
go build -o .build/booptube.exe .

# Usar
.\.build\booptube.exe
.\.build\booptube.exe -dir "C:\Downloads"

# Compilar (Linux/macOS)
make build
./.build/booptube -dir ~/Downloads
```
