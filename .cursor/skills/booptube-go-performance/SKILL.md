---
name: booptube-go-performance
description: >-
  Optimizes Go code for lower RAM use and faster runtime: escape analysis, slice/map
  preallocation, sync.Pool, goroutine lifecycle, buffered channels, bufio I/O, and
  GOMEMLIMIT/GOGC. Use when profiling, optimizing hot paths, reducing allocations,
  or reviewing Go performance in the booptube repository.
---

# Booptube Go Performance

Você é um agente que busca sempre melhorar a performance, deixando o software consumindo menos memória RAM e otimizando seu algoritmo para ser executado em runtime.

Priorize mudanças mensuráveis em caminhos quentes. Não micro-otimize código frio. Combine com `booptube-go-style` (estrutura) e `booptube-go-resilience` (estabilidade, timeouts, OOM em deploy).

## 1. Alocação de Memória e Gerenciamento de Heap vs Stack

O Garbage Collector (GC) do Go é altamente otimizado, mas quanto menos trabalho você der a ele, mais rápida será sua aplicação.

### Evite Escapes para o Heap (Escape Analysis)

Sempre que você usa ponteiros (`*`), corre o risco de mandar a variável para o Heap (memória compartilhada gerenciada pelo GC) em vez de mantê-la na Stack (memória local da função, que é virtualmente gratuita).

**Regra:** Passe estruturas pequenas por valor, não por ponteiro. Só use ponteiros se precisar modificar o estado original ou se a estrutura for muito grande (acima de alguns kilobytes).

```go
// ❌ Evitar — pequena struct no heap sem necessidade
func process(v *VideoMeta) { ... }

// ✅ Preferir — valor na stack
func process(v VideoMeta) { ... }
```

Verifique escapes com `go build -gcflags="-m" ./...` antes de refatorar tipos públicos.

### Pré-aloque Slices e Maps

Expandir um slice ou map dinamicamente força o Go a alocar um novo bloco de memória maior e copiar os dados antigos.

**Regra:** Se você já sabe o tamanho final (ou uma estimativa), use `make([]Tipo, 0, capacidade)` ou `make(map[Chave]Valor, capacidade)`.

```go
// ❌ Evitar
var ids []string
for _, v := range videos {
	ids = append(ids, v.ID)
}

// ✅ Preferir
ids := make([]string, 0, len(videos))
for _, v := range videos {
	ids = append(ids, v.ID)
}
```

### Reutilize Objetos com sync.Pool

Para objetos que são criados e destruídos milhares de vezes por segundo (como buffers de bytes em requisições HTTP), use um pool de objetos. Isso reduz drasticamente a fragmentação de memória e o estresse sobre o GC.

```go
var bufPool = sync.Pool{
	New: func() any { return make([]byte, 32*1024) },
}

func readChunk(r io.Reader) ([]byte, error) {
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)
	n, err := r.Read(buf)
	return buf[:n], err // só seguro se o caller copiar antes de devolver ao pool
}
```

Use `sync.Pool` apenas para buffers temporários; nunca guarde referências ao objeto após `Put`.

## 2. Concorrência Eficiente (Goroutines e Schedulers)

Goroutines são baratas (começam com cerca de 2 KB de memória), mas não são de graça.

### Evite o Vazamento de Goroutines (Goroutine Leaks)

Uma goroutine que fica travada esperando um canal que nunca será fechado ou um lock que nunca será liberado permanece na memória para sempre.

**Regra:** Sempre defina políticas de encerramento usando `context.WithTimeout` ou `context.WithCancel`. Certifique-se de que toda goroutine criada tenha uma condição clara de saída.

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()

go func() {
	select {
	case <-ctx.Done():
		return
	case result := <-workCh:
		handle(result)
	}
}()
```

### Cuidado com Canais Sem Buffer (Unbuffered Channels)

Canais sem buffer bloqueiam a goroutine emissora até que a receptora leia o dado. Se o processamento do receptor for lento, você cria um gargalo em cadeia.

**Regra:** Use canais com buffer (`make(chan Tipo, tamanho)`) quando puder prever o volume, permitindo que os emissores continuem trabalhando sem esperar o consumo imediato.

### Não abuse de runtime.GOMAXPROCS

Por padrão, o Go define o número de threads do sistema operacional igual ao número de núcleos de CPU disponíveis. Alterar isso manualmente para um valor muito alto causa context switching excessivo, degradando a performance.

Só ajuste `GOMAXPROCS` com benchmark (`go test -bench`) ou perfil de CPU que comprove ganho.

## 3. Otimização de I/O (Entrada e Saída)

Operações de disco e rede costumam ser os maiores gargalos de qualquer sistema.

### Use Buffers para Leitura e Escrita

Nunca leia ou escreva em arquivos ou conexões byte a byte diretamente.

**Regra:** Envolva seus leitores e escritores com o pacote `bufio` (ex: `bufio.NewReader` e `bufio.NewWriter`). Isso agrupa as operações em blocos maiores de memória, minimizando chamadas de sistema (syscalls).

```go
w := bufio.NewWriterSize(f, 256*1024)
defer w.Flush()
```

### Evite Conversões Inúteis de string para []byte

Converter uma string para um slice de bytes (ou vice-versa) força o Go a fazer uma nova alocação e copiar os dados na memória.

**Regra:** Se você está trabalhando com I/O de rede ou arquivos, tente manter os dados como `[]byte` do início ao fim. Use pacotes como `bytes` em vez de `strings` para manipulação.

```go
// ❌ Evitar em loop quente
body := []byte(fmt.Sprintf("id=%s", id))

// ✅ Preferir
var b bytes.Buffer
b.WriteString("id=")
b.WriteString(id)
```

## 4. Uso de Variáveis Globais de Controle (A partir do Go 1.19+)

O runtime moderno do Go introduziu ferramentas para blindar a performance a nível de infraestrutura:

**GOMEMLIMIT:** Define um teto rígido de memória para o runtime do Go. Se a sua aplicação estiver próxima desse limite, o Go executará o Garbage Collector de forma muito mais agressiva para evitar que o sistema operacional mate o seu processo por falta de memória (OOM Kill).

**GOGC:** Controla a agressividade do GC baseado no crescimento percentual da memória alocada. O padrão é 100. Se você mudar para `off`, desativa o GC (útil apenas para scripts de execução curta); se definir valores menores (ex: 50), o GC roda mais vezes, economizando memória ao custo de um pouco mais de CPU.

Documente valores de deploy em `config/` ou README (ex: `GOMEMLIMIT=1800MiB`, `GOGC=100`).

## Tabela de Boas Práticas Rápidas

| O que EVITAR ❌ | O que FAZER CORRETAMENTE ✅ |
|-----------------|----------------------------|
| `append(slice, item)` em loops longos | `slice := make([]T, 0, len)`; `append(...)` |
| `fmt.Sprintf("id: %d", id)` em loops internos | `strconv.Itoa(id)` ou `strconv.AppendInt()` |
| Abrir arquivos/conexões sem `defer close` | Sempre usar `defer file.Close()` imediatamente |
| Usar reflexão (`reflect`) em caminhos críticos | Código fortemente tipado ou geração de código |

## Workflow de otimização

1. **Medir primeiro:** `go test -bench`, `pprof` (CPU + heap), ou `go build -gcflags="-m"`.
2. **Atacar o maior custo:** alocações em loop > I/O sem buffer > goroutines presas > micro-refactors.
3. **Validar:** benchmark antes/depois; heap profile não deve piorar sem ganho de CPU claro.
4. **Entregar diff mínimo:** uma otimização por PR quando possível.

## Checklist antes de entregar

- [ ] Caminhos quentes pré-alocam slices/maps com capacidade conhecida
- [ ] Structs pequenas passadas por valor; ponteiros só quando necessário
- [ ] Buffers HTTP/disco usam `bufio` ou `io.CopyBuffer` com buffer reutilizável
- [ ] Goroutines têm saída via `context` ou fechamento de canal
- [ ] Canais bufferizados quando produtor e consumidor têm ritmos diferentes
- [ ] Loops internos evitam `fmt.Sprintf` e conversões string↔[]byte desnecessárias
- [ ] Sem `reflect` em hot path
- [ ] `defer Close()` imediato após abrir arquivo/conexão
- [ ] Deploy documenta `GOMEMLIMIT` / `GOGC` quando relevante

## Recursos adicionais

- Padrões de buffer, pool e benchmark: [examples.md](examples.md)
