# Guia do usuário — booptube

O **booptube** baixa vídeos do YouTube. Você escolhe a pasta de destino, cola o link e escolhe **mp4** (vídeo) ou **mp3** (só áudio).

Há duas formas de usar:

| Modo | Executável | Descrição |
|------|------------|-----------|
| **Terminal (CLI)** | `booptube` / `booptube.exe` | Prompts interativos no terminal |
| **Interface gráfica (GUI)** | `booptube-gui` / `booptube-gui.exe` | Janela com visual neon futurista |

Não é preciso instalar ffmpeg nem yt-dlp — tudo já vem dentro do programa.

---

## Antes de começar

Você precisa apenas do executável:

| Sistema | CLI | GUI |
|---------|-----|-----|
| Windows | `booptube.exe` | `booptube-gui.exe` |
| Linux / macOS | `booptube` | `booptube-gui` |

Abra um terminal (PowerShell, CMD, Terminal) na pasta onde está o executável, ou use o caminho completo. A GUI abre com duplo clique no explorador de arquivos.

**Primeira execução:** o programa pode demorar alguns segundos enquanto prepara os componentes internos. Isso é normal e acontece só uma vez (ou após atualizar o booptube).

---

## Interface gráfica (GUI)

Abra **`booptube-gui`** (ou `booptube-gui.exe` no Windows). A janela tem fundo escuro com acentos neon.

### Passo a passo na GUI

1. **Pasta de destino** — digite o caminho ou clique em **Escolher...** para selecionar uma pasta. A última pasta usada é carregada automaticamente.
2. **URL do YouTube** — cole o link do vídeo.
3. **Formato** — selecione **MP4 (vídeo)** ou **MP3 (áudio)**.
4. Clique em **Baixar**.

Durante o download:

- A barra de progresso mostra o percentual quando disponível.
- O **Log** exibe as mensagens do yt-dlp.
- Use **Cancelar** para interromper.

Ao concluir, a mensagem **Concluído.** aparece e você pode baixar outro vídeo sem fechar a janela.

### Encerrar a GUI

Feche a janela normalmente. Se houver download em andamento, ele será cancelado.

---

## Como abrir o booptube (CLI)

### Windows (PowerShell ou CMD)

```powershell
.\booptube.exe
```

### Linux / macOS

```bash
./booptube
```

Se aparecer a mensagem inicial:

```text
booptube — digite q ou sair para encerrar
```

O programa está pronto. Siga os passos abaixo.

---

## Passo a passo

A cada download, o booptube faz **3 perguntas** (ou **2**, se você usar a flag `-dir`).

### 1. Pasta de destino

```text
Pasta de destino:
```

Digite o caminho da pasta onde os arquivos serão salvos.

**Exemplos (Windows):**

```text
C:\Users\voce\Downloads
D:\Videos\YouTube
```

**Exemplos (Linux / macOS):**

```text
/home/voce/Downloads
/Users/voce/Downloads
```

- A pasta é **criada automaticamente** se não existir.
- Se você já baixou algo antes, verá `(Enter=C:\caminho\anterior)` — pressione **Enter** para usar a mesma pasta.
- A última pasta usada fica salva para a próxima vez que abrir o programa.

### 2. URL do YouTube

```text
URL do YouTube:
```

Cole o link do vídeo. Exemplos válidos:

```text
https://www.youtube.com/watch?v=VIDEO_ID
https://youtu.be/VIDEO_ID
www.youtube.com/watch?v=VIDEO_ID
```

**Aceito:** links de `youtube.com`, `youtu.be`, YouTube Music e links curtos.

**Não aceito:** sites que não sejam YouTube, playlists inteiras (só o vídeo do link é baixado).

### 3. Formato

```text
Formato [1=mp4, 2=mp3] (Enter=mp4):
```

| O que digitar | Resultado |
|---------------|-----------|
| **Enter** (vazio) | Vídeo em **mp4** |
| `1` ou `mp4` | Vídeo em **mp4** |
| `2` ou `mp3` | Áudio em **mp3** |

Depois disso o download começa. O progresso aparece no terminal. Ao terminar:

```text
concluido.
```

O programa **pergunta de novo** a pasta, URL e formato — você pode baixar outro vídeo sem fechar.

---

## Exemplo completo

```text
booptube — digite q ou sair para encerrar
Pasta de destino: C:\Downloads\booptube
URL do YouTube: https://www.youtube.com/watch?v=jNQXAC9IVRw
Formato [1=mp4, 2=mp3] (Enter=mp4): 2
baixando https://www.youtube.com/watch?v=jNQXAC9IVRw como mp3...
[download] 100% ...
concluido.
Pasta de destino (Enter=C:\Downloads\booptube):
```

Arquivo gerado: `C:\Downloads\booptube\Me at the zoo.mp3`

O nome do arquivo é o **título do vídeo** no YouTube.

---

## Opção: definir a pasta ao abrir

Se você sempre salva na mesma pasta, use a flag **`-dir`** ao iniciar. O booptube **não pergunta** a pasta de destino.

### Windows

```powershell
.\booptube.exe -dir "C:\Downloads\booptube"
```

### Linux / macOS

```bash
./booptube -dir "$HOME/Downloads/booptube"
```

Depois disso só aparecem as perguntas de **URL** e **formato**.

---

## Como sair

| Ação | Quando usar |
|------|-------------|
| Digite `q`, `sair`, `quit` ou `exit` | Em qualquer prompt de texto |
| `Ctrl+C` | Para cancelar um download ou fechar o programa |

---

## Onde ficam os arquivos

### Downloads

Na pasta que você escolheu (ou passou com `-dir`).

| Formato | Exemplo de nome |
|---------|-----------------|
| mp4 | `Nome do video.mp4` |
| mp3 | `Nome do video.mp3` |

### Configuração (pasta lembrada)

O booptube guarda a última pasta de destino automaticamente:

| Sistema | Arquivo |
|---------|---------|
| Windows | `%AppData%\booptube\config.json` |
| Linux / macOS | `~/.config/booptube/config.json` |

Conteúdo típico:

```json
{"download_dir":"C:\\Users\\voce\\Downloads\\booptube"}
```

Você **não precisa editar** esse arquivo — a CLI atualiza sozinha. Se quiser fixar uma pasta padrão manualmente, edite `download_dir` antes de abrir o booptube.

---

## Comandos disponíveis

### Ao iniciar o programa

```text
booptube [-dir pasta]
```

| Opção | Descrição |
|-------|-----------|
| `-dir pasta` | Define a pasta de destino e pula essa pergunta |
| `-h` | Mostra ajuda breve |

### Durante o uso (nos prompts)

| Entrada | Efeito |
|---------|--------|
| Caminho de pasta | Define onde salvar |
| Enter (com pasta anterior) | Reutiliza a última pasta |
| URL do YouTube | Inicia o download daquele vídeo |
| `1`, `mp4`, Enter | Baixa vídeo |
| `2`, `mp3` | Baixa só áudio |
| `q`, `sair`, `quit`, `exit` | Fecha o booptube |

---

## Problemas comuns

### `url vazia` ou `url invalida`

Cole um link completo de vídeo do YouTube. Evite links de playlist ou páginas que não sejam vídeo.

### `apenas URLs do YouTube sao suportadas`

O link não é do YouTube. Use `youtube.com` ou `youtu.be`.

### `pasta nao gravavel`

Escolha outra pasta (por exemplo `Downloads`) ou verifique permissões de escrita no disco.

### `download falhou`

Possíveis causas:

- Vídeo privado, removido ou com restrição regional
- Sem internet ou conexão instável
- Link incorreto

Tente outro vídeo ou confira se o link abre no navegador.

### Download muito lento

Depende da sua internet e do tamanho do vídeo. O progresso no terminal mostra o andamento — aguarde ou cancele com `Ctrl+C`.

### Quero mudar a pasta padrão

Na próxima execução, digite a nova pasta quando perguntado (ou use `-dir "nova\pasta"`). A configuração será atualizada após o próximo download concluído.

### GUI não abre ou não compila

Veja o guia dedicado: **[gui.md](gui.md)** — instalação do MinGW (Windows), compilação com CGO e problemas comuns.

---

## Dicas rápidas

1. Use **`-dir`** se sempre baixa para o mesmo lugar.
2. Pressione **Enter** na pasta para repetir o último destino.
3. Pressione **Enter** no formato para baixar **mp4** (padrão).
4. Baixe **um vídeo por vez** — cole o link de um vídeo, não de playlist.
5. O arquivo final usa o **título do vídeo**; caracteres inválidos no Windows são ajustados automaticamente pelo yt-dlp.

---

## Precisa compilar ou desenvolver?

Este guia é só para **usar** o programa.

| Documento | Conteúdo |
|-----------|----------|
| [gui.md](gui.md) | Instalar e compilar a GUI |
| [projeto.md](projeto.md) | Visão geral completa do projeto |
| [cli.md](cli.md) | Build, Makefile e detalhes técnicos |
