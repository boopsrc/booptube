package web

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"booptube/config"
	"booptube/downloader"
	"booptube/video"
)

type Server struct {
	cfg    config.WebConfig
	dl     *downloader.Client
	mu     sync.Mutex
	jobs   map[string]*Job
	queue  []string
	sem    chan struct{}
	stopCh chan struct{}
}

type Job struct {
	ID            string    `json:"id"`
	URL           string    `json:"url"`
	Format        string    `json:"format"`
	Status        string    `json:"status"`
	Progress      float64   `json:"progress"`
	FilePath      string    `json:"-"`
	Filename      string    `json:"filename,omitempty"`
	Error         string    `json:"error,omitempty"`
	Log           []string  `json:"log,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	FileCreatedAt time.Time `json:"file_created_at,omitempty"`
	ReadyAt       time.Time `json:"ready_at,omitempty"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
	TTLSeconds    int       `json:"ttl_seconds"`
	PageURL       string    `json:"page_url,omitempty"`
	DownloadURL   string    `json:"download_url,omitempty"`
}

type createRequest struct {
	URL    string `json:"url"`
	Format string `json:"format"`
}

func New(cfg config.WebConfig, dl *downloader.Client) *Server {
	s := &Server{
		cfg:    cfg,
		dl:     dl,
		jobs:   make(map[string]*Job),
		sem:    make(chan struct{}, cfg.MaxConcurrent),
		stopCh: make(chan struct{}),
	}
	go s.workerLoop()
	go s.sweeperLoop()
	return s
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleIndex)
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS()))))
	mux.HandleFunc("POST /api/downloads", s.handleCreate)
	mux.HandleFunc("GET /api/downloads/{id}", s.handleStatus)
	mux.HandleFunc("GET /api/downloads/{id}/file", s.handleFile)
	mux.HandleFunc("GET /d/{id}", s.handleDownloadPage)
	mux.HandleFunc("GET /d/{id}/file", s.handleFile)
	return s.logMiddleware(mux)
}

func (s *Server) Shutdown() {
	close(s.stopCh)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	serveStaticHTML(w, "static/index.html")
}

func (s *Server) handleDownloadPage(w http.ResponseWriter, r *http.Request) {
	serveStaticHTML(w, "static/download.html")
}

func serveStaticHTML(w http.ResponseWriter, name string) {
	data, err := staticFS.ReadFile(name)
	if err != nil {
		http.Error(w, "pagina indisponivel", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json invalido"})
		return
	}
	parsedURL, err := video.ParseURL(req.URL)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	format, err := video.FormatFromString(req.Format)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := newJobID()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "falha ao criar job"})
		return
	}

	pageURL := "/d/" + id
	job := &Job{
		ID:         id,
		URL:        parsedURL,
		Format:     format.String(),
		Status:     "queued",
		CreatedAt:  time.Now().UTC(),
		PageURL:    pageURL,
		TTLSeconds: int(s.cfg.FileTTL.Seconds()),
	}

	s.mu.Lock()
	s.jobs[id] = job
	s.queue = append(s.queue, id)
	s.mu.Unlock()

	slog.Info("download_queued",
		"job_id", id,
		"url", parsedURL,
		"format", format.String(),
		"client_ip", clientIP(r),
	)

	writeJSON(w, http.StatusAccepted, map[string]any{
		"id":       id,
		"status":   "queued",
		"page_url": pageURL,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job := s.getJob(id)
	if job == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "job nao encontrado"})
		return
	}
	writeJSON(w, http.StatusOK, job.snapshot())
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.serveJobFile(w, r, id)
}

func (s *Server) serveJobFile(w http.ResponseWriter, r *http.Request, id string) {
	job := s.getJob(id)
	if job == nil {
		http.Error(w, "job nao encontrado", http.StatusNotFound)
		return
	}

	s.mu.Lock()
	status := job.Status
	filePath := job.FilePath
	filename := job.Filename
	expiresAt := job.ExpiresAt
	s.mu.Unlock()

	if status != "ready" {
		if status == "expired" {
			http.Error(w, "arquivo expirado", http.StatusGone)
			return
		}
		http.Error(w, "arquivo ainda nao disponivel", http.StatusConflict)
		return
	}
	if time.Now().After(expiresAt) {
		s.expireJob(id)
		http.Error(w, "arquivo expirado", http.StatusGone)
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "arquivo indisponivel", http.StatusGone)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Type", contentTypeFor(filename))
	slog.Info("file_served",
		"job_id", id,
		"filename", filename,
		"client_ip", clientIP(r),
	)
	io.Copy(w, f)
}

func (s *Server) workerLoop() {
	for {
		select {
		case <-s.stopCh:
			return
		default:
		}

		id := s.dequeue()
		if id == "" {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		s.sem <- struct{}{}
		safeGo(func() {
			defer func() { <-s.sem }()
			s.runJob(id)
		})
	}
}

func (s *Server) runJob(id string) {
	s.mu.Lock()
	job, ok := s.jobs[id]
	if !ok {
		s.mu.Unlock()
		return
	}
	job.Status = "downloading"
	jobDir := filepath.Join(s.cfg.DownloadDir, id)
	url := job.URL
	formatStr := job.Format
	s.mu.Unlock()

	if err := os.MkdirAll(jobDir, 0o755); err != nil {
		s.failJob(id, fmt.Errorf("criar pasta do job: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.DownloadTimeout)
	defer cancel()

	start := time.Now()
	slog.Info("download_started", "job_id", id, "url", url, "format", formatStr)

	var logMu sync.Mutex
	var logLines []string
	updateJob := func(fn func(*Job)) {
		s.mu.Lock()
		defer s.mu.Unlock()
		if j := s.jobs[id]; j != nil {
			fn(j)
		}
	}
	handlers := &downloader.Handlers{
		OnLine: func(line string) {
			logMu.Lock()
			logLines = append(logLines, line)
			if len(logLines) > 20 {
				logLines = logLines[len(logLines)-20:]
			}
			lines := append([]string(nil), logLines...)
			logMu.Unlock()

			updateJob(func(j *Job) {
				j.Log = lines
				if pct, ok := downloader.ParseProgress(line); ok {
					j.Progress = pct
				}
			})
		},
		OnProgress: func(pct float64) {
			logMu.Lock()
			lines := append([]string(nil), logLines...)
			logMu.Unlock()
			updateJob(func(j *Job) {
				j.Progress = pct
				j.Log = lines
			})
		},
	}

	format, _ := video.FormatFromString(formatStr)
	outFile := filepath.Join(jobDir, id+"."+formatStr)
	filePath, err := s.dl.DownloadToFile(ctx, url, format, outFile, handlers)
	if err != nil {
		s.failJob(id, err)
		os.RemoveAll(jobDir)
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		s.failJob(id, fmt.Errorf("stat arquivo: %w", err))
		os.RemoveAll(jobDir)
		return
	}

	fileCreatedAt := info.ModTime().UTC()
	expiresAt := fileCreatedAt.Add(s.cfg.FileTTL)
	filename := filepath.Base(filePath)
	pageURL := "/d/" + id
	downloadURL := pageURL + "/file"

	updateJob(func(j *Job) {
		j.Status = "ready"
		j.Progress = 100
		j.FilePath = filePath
		j.Filename = filename
		j.ReadyAt = time.Now().UTC()
		j.FileCreatedAt = fileCreatedAt
		j.ExpiresAt = expiresAt
		j.PageURL = pageURL
		j.DownloadURL = downloadURL
		logMu.Lock()
		j.Log = append([]string(nil), logLines...)
		logMu.Unlock()
	})

	slog.Info("download_ready",
		"job_id", id,
		"filename", filename,
		"duration_ms", time.Since(start).Milliseconds(),
		"file_created_at", fileCreatedAt.Format(time.RFC3339),
		"expires_at", expiresAt.Format(time.RFC3339),
	)

	time.AfterFunc(time.Until(expiresAt), func() {
		s.expireJob(id)
	})
}

func (s *Server) failJob(id string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j := s.jobs[id]; j != nil {
		j.Status = "failed"
		j.Error = err.Error()
	}
	slog.Error("download_failed", "job_id", id, "error", err.Error())
}

func (s *Server) expireJob(id string) {
	s.mu.Lock()
	job, ok := s.jobs[id]
	if !ok || job.Status == "expired" {
		s.mu.Unlock()
		return
	}
	filePath := job.FilePath
	jobDir := filepath.Join(s.cfg.DownloadDir, id)
	job.Status = "expired"
	job.FilePath = ""
	job.DownloadURL = ""
	s.mu.Unlock()

	if filePath != "" {
		os.Remove(filePath)
	}
	os.RemoveAll(jobDir)
	slog.Info("job_expired", "job_id", id)
}

func (s *Server) sweeperLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case now := <-ticker.C:
			s.sweep(now)
		}
	}
}

func (s *Server) sweep(now time.Time) {
	s.mu.Lock()
	var expired []string
	for id, job := range s.jobs {
		if job.Status == "ready" && !job.ExpiresAt.IsZero() && now.After(job.ExpiresAt) {
			expired = append(expired, id)
		}
	}
	s.mu.Unlock()

	for _, id := range expired {
		s.expireJob(id)
	}

	entries, err := os.ReadDir(s.cfg.DownloadDir)
	if err != nil {
		return
	}
	s.mu.Lock()
	known := make(map[string]bool, len(s.jobs))
	for id := range s.jobs {
		known[id] = true
	}
	s.mu.Unlock()

	for _, e := range entries {
		if !e.IsDir() || known[e.Name()] {
			continue
		}
		orphan := filepath.Join(s.cfg.DownloadDir, e.Name())
		os.RemoveAll(orphan)
		slog.Info("job_cleaned", "path", orphan, "reason", "orphan")
	}
}

func (s *Server) dequeue() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.queue) == 0 {
		return ""
	}
	id := s.queue[0]
	s.queue = s.queue[1:]
	return id
}

func (s *Server) getJob(id string) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok {
		return nil
	}
	cp := *j
	cp.Log = append([]string(nil), j.Log...)
	return &cp
}

func (j *Job) snapshot() Job {
	cp := *j
	cp.Log = append([]string(nil), j.Log...)
	return cp
}

func newJobID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func contentTypeFor(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	default:
		return "application/octet-stream"
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i >= 0 {
		return host[:i]
	}
	return host
}

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", clientIP(r),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func safeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("goroutine_panic", "panic", fmt.Sprint(r))
			}
		}()
		fn()
	}()
}
