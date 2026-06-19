# Exemplos — Booptube Go Resilience

## Download com semáforo, timeout e recover

Padrão recomendado para `downloader/`:

```go
type Client struct {
	sem    chan struct{}
	client *http.Client
}

func NewClient(maxConcurrent int) *Client {
	return &Client{
		sem:    make(chan struct{}, maxConcurrent),
		client: &http.Client{Timeout: 0}, // deadline vem do context
	}
}

func (c *Client) Download(ctx context.Context, url, dest string) (err error) {
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return ctx.Err()
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmp, err := os.CreateTemp(filepath.Dir(dest), ".download-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		tmp.Close()
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	if _, err = io.Copy(tmp, resp.Body); err != nil {
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}

	return withRetry(5, func() error {
		return os.Rename(tmpPath, dest)
	})
}

func (c *Client) DownloadAsync(parent context.Context, url, dest string, onDone func(error)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				onDone(fmt.Errorf("panic in download: %v", r))
				return
			}
		}()
		onDone(c.Download(parent, url, dest))
	}()
}
```

## Bootstrap em main.go

```go
func main() {
	if runtime.GOOS == "darwin" {
		signal.Notify(make(chan os.Signal, 1), syscall.SIGPIPE)
	}

	// wiring...
}
```
