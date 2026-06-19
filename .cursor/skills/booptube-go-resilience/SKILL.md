---
name: booptube-go-resilience
description: >-
  Enforces Go resilience and cross-platform performance for booptube: goroutine
  recover, I/O timeouts, GOMEMLIMIT, file-descriptor pools, Windows file-lock
  retries, CGO LockOSThread, macOS App Nap and SIGPIPE handling. Use when
  writing, reviewing, or debugging Go code for crashes, panics, OOM kills, too
  many open files, Windows locking, or macOS network/process issues.
---

# Booptube Go Resilience

Para garantir que um software em Go (Golang) seja resiliente a falhas (nunca crashe) e mantenha alta performance rodando nativamente em Windows, Linux e macOS, precisamos blindar o código contra os comportamentos específicos de cada núcleo de sistema operacional.

Em Go, o runtime gerencia threads (através de goroutines), memória (via Garbage Collector) e chamadas de sistema. Se esses recursos não forem controlados, o OS vai encerrar o seu processo à força (gerando os temidos panics ou OOM Kills).

Abaixo estão as regras fundamentais de performance e estabilidade divididas por ecossistema e arquitetura.

## Regra de Ouro Universal (Cross-Platform)

Antes das especificidades de cada OS, seu código Go deve seguir esta lei para evitar panics:

**Recuperação de Goroutines:** Qualquer nova goroutine (`go func()`) deve ter um bloco `defer recover()` interno. Um panic não capturado dentro de uma goroutine isolada derruba a aplicação inteira.

**Gargalos de I/O:** Sempre utilize `context.WithTimeout` para operações de rede ou leitura de disco. Travamentos de I/O consomem threads do sistema operacional (`M` no modelo MPG do Go), podendo esgotar os limites do OS.

### Padrões obrigatórios

```go
func safeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("goroutine panic: %v", r)
			}
		}()
		fn()
	}()
}
```

```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
if err != nil {
	return err
}
resp, err := http.DefaultClient.Do(req)
if err != nil {
	return err
}
defer resp.Body.Close()
```

## 1. Regras para Linux (Foco em Gerenciamento de Recursos)

O Linux é extremamente rigoroso com o uso de memória e descritores de arquivo. Se o seu app abusar, o kernel vai matá-lo sem aviso.

**Evite o OOM Killer (Out of Memory):** O Linux encerra processos que consomem memória de forma rampa e descontrolada.

**Ação:** Defina a variável de ambiente `GOMEMLIMIT` (disponível a partir do Go 1.19). Se o seu container ou máquina tem 2GB de RAM, configure `GOMEMLIMIT=1800MiB`. Isso força o Garbage Collector do Go a trabalhar agressivamente antes que o Linux decida matar o processo.

**Estouro de File Descriptors (Limits):** No Linux, "tudo é um arquivo" (sockets de rede, arquivos de vídeo, etc.). O limite padrão por processo costuma ser baixo (1024). Se você baixar múltiplos vídeos simultaneamente, o app vai falhar com `too many open files`.

**Ação:** Use pools de conexões (como um `chan struct{}` atuando como semáforo) para limitar downloads simultâneos e sempre feche os corpos de requisições com `defer resp.Body.Close()`.

### Padrões Linux

```go
type Semaphore struct {
	slots chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{slots: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire() { s.slots <- struct{}{} }
func (s *Semaphore) Release() { <-s.slots }

// Limitar downloads simultâneos (ex: 8)
sem := NewSemaphore(8)
sem.Acquire()
defer sem.Release()
```

Documentar em `config/` ou README: `GOMEMLIMIT=1800MiB` para deploy Linux/container.

## 2. Regras para Windows (Foco em File Locking e Syscalls)

O subsistema de arquivos do Windows (NTFS) e o gerenciamento de processos funcionam de forma muito diferente do modelo POSIX (Linux/Mac).

**Tratamento de File Locking (Arquivos Presos):** No Windows, se o seu programa estiver escrevendo em um arquivo de vídeo e outra rotina tentar editá-lo ou movê-lo, o Windows bloqueará a operação e retornará um erro de permissão que pode quebrar o app.

**Ação:** Implemente retries com exponential backoff para operações de escrita/leitura de arquivos e garanta que os ponteiros de arquivo (`os.File`) sejam fechados imediatamente após o uso, nunca acumulados em memória.

**Gerenciamento de Threads com CGO:** Se o seu design system futurista interagir com APIs nativas do Windows (como Win32 ou DirectX para efeitos visuais) via CGO, essas chamadas rodam em threads do OS presas (locked).

**Ação:** Use `runtime.LockOSThread()` na goroutine que interage com a interface gráfica do Windows para garantir que o Go não mova a execução para outra thread, o que causaria crashes de violação de acesso à memória (Access Violation).

### Padrões Windows

```go
func withRetry(max int, fn func() error) error {
	var err error
	backoff := 50 * time.Millisecond
	for i := 0; i < max; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return err
}

// Uso: mover/renomear arquivo após download
err := withRetry(5, func() error {
	return os.Rename(tmpPath, finalPath)
})
```

```go
// Apenas na goroutine que chama Win32/CGO/DirectX
runtime.LockOSThread()
defer runtime.UnlockOSThread()
```

Sempre `defer f.Close()` logo após `os.Open` / `os.Create`; nunca guardar `*os.File` em structs de longa duração.

## 3. Regras para macOS / Darwin (Foco em Sleeping e Arquitetura ARM)

O macOS possui regras estritas de economia de energia e transição de arquiteturas (Intel vs Apple Silicon).

**Prevenção do App Nap (Processo Congelado):** O macOS coloca aplicações de background ou de terminal inativas em modo "App Nap", reduzindo drasticamente os ciclos de CPU. Se o seu downloader for acordado abruptamente com uma rajada de dados, o desalinhamento de timers pode quebrar conexões de rede.

**Ação:** Para ferramentas CLI ou de background que realizam tarefas longas, utilize syscalls específicas (ou bibliotecas que envelopam a API do macOS) para sinalizar que o processo está realizando uma atividade de alta prioridade (Activity User Initiated).

**Sinais de Sistema (SIGPIPE):** No macOS, tentar escrever em uma conexão de rede ou pipe que já foi fechada pelo servidor gera um sinal SIGPIPE, que por padrão finaliza o programa imediatamente.

**Ação:** O runtime do Go lida com isso na maioria das vezes, mas ao usar pacotes customizados de rede ou CGO, certifique-se de ignorar ou tratar adequadamente os sinais do sistema operacional usando o pacote `os/signal`.

### Padrões macOS

```go
// main.go — ignorar SIGPIPE quando usar CGO ou sockets raw
signal.Notify(make(chan os.Signal, 1), syscall.SIGPIPE)
```

Para downloads longos em CLI macOS, considerar wrapper de Activity (ex: `github.com/prashantgupta17/mac-sleep-prevent`) ou documentar que o operador deve desabilitar App Nap para sessões batch críticas.

## Checklist antes de entregar

- [ ] Toda `go func()` tem `defer recover()` (direto ou via helper `safeGo`)
- [ ] HTTP/disco/arquivo usam `context.WithTimeout` ou deadline equivalente
- [ ] Downloads simultâneos limitados por semáforo; `resp.Body.Close()` sempre com `defer`
- [ ] Deploy Linux documenta `GOMEMLIMIT` proporcional à RAM disponível
- [ ] Operações de arquivo no Windows usam retry com backoff; `os.File` fechado no escopo
- [ ] CGO/Win32/UI nativa usa `runtime.LockOSThread()` na goroutine correta
- [ ] Rede customizada/CGO no macOS trata `SIGPIPE` via `os/signal`
- [ ] Nenhum panic não recuperado pode derrubar o processo inteiro

## Recursos adicionais

- Exemplos completos de download com semáforo + timeout: [examples.md](examples.md)
