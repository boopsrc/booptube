# Exemplos — Booptube Go Performance

## Download com bufio, buffer reutilizável e pré-alocação

Padrão recomendado para `downloader/`:

```go
const copyBufSize = 256 * 1024

var copyBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, copyBufSize)
		return &b
	},
}

func copyToFile(dst *os.File, src io.Reader) (written int64, err error) {
	p := copyBufPool.Get().(*[]byte)
	buf := *p
	defer copyBufPool.Put(p)

	w := bufio.NewWriterSize(dst, copyBufSize)
	defer func() {
		if flushErr := w.Flush(); err == nil {
			err = flushErr
		}
	}()

	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			if _, err = w.Write(buf[:n]); err != nil {
				return written, err
			}
			written += int64(n)
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return written, readErr
		}
	}
	return written, nil
}
```

## Worker pool com canal bufferizado e context

```go
func processAll(ctx context.Context, urls []string, workers int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan string, len(urls))
	errs := make(chan error, workers)

	for w := 0; w < workers; w++ {
		go func() {
			for url := range jobs {
				if err := downloadOne(ctx, url); err != nil {
					select {
					case errs <- err:
					default:
					}
					cancel()
					return
				}
			}
		}()
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

## Benchmark antes/depois

```go
func BenchmarkFormatID(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = strconv.FormatInt(int64(i), 10)
	}
}
```

Rodar: `go test -bench=. -benchmem ./downloader/`
