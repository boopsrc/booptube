# booptube — documentação da CLI

CLI interativa para baixar vídeos do YouTube em **mp4** ou **mp3**. O binário embute o **yt-dlp**; na primeira execução ele é extraído para o cache do usuário.

## Índice

1. [Como funciona](#como-funciona)
2. [Instalação](#instalação)
3. [Configuração](#configuração)
4. [Uso](#uso)
5. [Comandos e flags](#comandos-e-flags)
6. [Prompts interativos](#prompts-interativos)
7. [Comandos do Makefile](#comandos-do-makefile)
8. [Arquivos gerados](#arquivos-gerados)
9. [Dependências](#dependências)
10. [Solução de problemas](#solução-de-problemas)

---

## Como funciona

```text
booptube inicia
    │
    ├─ Carrega config.json (última pasta, se existir)
    ├─ Extrai yt-dlp embutido para o cache (se necessário)
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
2. O booptube chama o **yt-dlp** (já incluído no executável) com os parâmetros corretos.
3. O arquivo é gravado na pasta escolhida com o nome `Título do vídeo.ext`.
4. Após cada download bem-sucedido, a **última pasta** é memorizada para a próxima sessão.

**Cancelamento:** `Ctrl+C` interrompe o download em andamento. Nos prompts, digite `q`, `sair`, `quit` ou `exit` para encerrar o programa.

**Playlists:** apenas o vídeo da URL informada é baixado (`--no-playlist`).

---

## Instalação

### Pré-requisitos

| Requisito | Versão mínima | Observação |
|-----------|---------------|------------|
| Go | 1.22+ | só para compilar a partir do código |
| ffmpeg | qualquer recente | **obrigatório** para mp3 e merge mp4 |
| make / curl | Linux e macOS | para `make build` |
| PowerShell | 5+ | para `fetch-ytdlp.ps1` no Windows |

### Instalar ffmpeg

**Windows (winget):**

```powershell
winget install Gyan.FFmpeg
```

Reabra o terminal após instalar para o `ffmpeg` entrar no PATH.

**Linux (Debian/Ubuntu):**

```bash
sudo apt install ffmpeg
```

**macOS (Homebrew):**

```bash
brew install ffmpeg
```

### Compilar a partir do código

#### Windows

```powershell
cd booptube
.\fetch-ytdlp.ps1
go build -o .build/booptube.exe .
```

O executável fica em `.build/booptube.exe`.

#### Linux / macOS

```bash
cd booptube
make build
```

Equivalente a `make fetch-ytdlp` + `go build -o .build/booptube .`.

O executável fica em `.build/booptube`.

> **Nota:** os binários do yt-dlp em `assets/ytdlp/` não vão para o git. É preciso rodar `fetch-ytdlp.ps1` ou `make fetch-ytdlp` antes do build.

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

### Variável de ambiente (build)

| Variável | Uso |
|----------|-----|
| `YTDLP_VERSION` | Versão do yt-dlp ao rodar `fetch-ytdlp.ps1` ou `make fetch-ytdlp` (padrão: `2026.06.09`) |

Exemplo:

```powershell
$env:YTDLP_VERSION = "2026.06.09"
.\fetch-ytdlp.ps1
```

---

## Uso

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

**mp4:** melhor vídeo + áudio, merge em mp4 (exige ffmpeg).

**mp3:** extrai áudio e converte para mp3 (exige ffmpeg).

---

## Comandos do Makefile

Disponíveis em Linux e macOS (requer `make` e `curl`):

| Comando | Descrição |
|---------|-----------|
| `make fetch-ytdlp` | Baixa yt-dlp para `assets/ytdlp/` (Windows, Linux e macOS arm64) |
| `make build` | Roda `fetch-ytdlp` e compila para `.build/booptube` |
| `make clean` | Remove a pasta `.build/` |

Versão pinada do yt-dlp no Makefile: `YTDLP_VERSION ?= 2026.06.09`.

Override:

```bash
make build YTDLP_VERSION=2026.06.09
```

### Script PowerShell (Windows)

| Script | Descrição |
|--------|-----------|
| `.\fetch-ytdlp.ps1` | Equivalente ao `make fetch-ytdlp` |

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

---

## Dependências

| Componente | Incluído no booptube? | Necessário? |
|------------|----------------------|-------------|
| yt-dlp | Sim (embutido) | Automático na extração |
| ffmpeg | Não | **Sim** — mp3 e merge mp4 |
| Go | Não | Só para compilar |

Sem ffmpeg, o download pode iniciar mas falha na conversão ou merge com mensagem do yt-dlp sobre `ffmpeg` / `ffprobe`.

---

## Solução de problemas

### `yt-dlp embutido ausente: execute fetch-ytdlp...`

O build foi feito sem baixar os assets. Rode `fetch-ytdlp.ps1` ou `make fetch-ytdlp` e recompile.

### `apenas URLs do YouTube sao suportadas`

A URL não é de um host YouTube reconhecido. Use link direto de vídeo (`watch?v=` ou `youtu.be/`).

### `pasta nao gravavel`

Escolha outra pasta ou corrija permissões. O booptube testa escrita antes de baixar.

### `Postprocessing: ffprobe and ffmpeg not found`

Instale ffmpeg e garanta que está no PATH do mesmo terminal onde roda o booptube.

### Download lento ou falha de rede

- Verifique conexão.
- URLs muito longas ou vídeos restritos podem falhar no yt-dlp.
- Recompile com yt-dlp mais recente: `make fetch-ytdlp` + `make build`.

### Atualizar yt-dlp embutido

1. Ajuste `YTDLP_VERSION` em `Makefile` / `fetch-ytdlp.ps1` se quiser pinar outra release.
2. Rode o fetch e recompile.
3. Na próxima execução, o cache local do yt-dlp será substituído.

---

## Referência rápida

```powershell
# Compilar (Windows)
.\fetch-ytdlp.ps1
go build -o .build/booptube.exe .

# Usar
.\.build\booptube.exe
.\.build\booptube.exe -dir "C:\Downloads"

# Compilar (Linux/macOS)
make build
./.build/booptube -dir ~/Downloads
```
