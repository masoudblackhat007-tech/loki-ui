package httpserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

const shutdownTimeout = 10 * time.Second

func Start(ctx context.Context, addr string) error {
	h := NewHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/readyz", h.Readyz)

	// UI
	mux.HandleFunc("/logs", h.LogsPage)
	mux.HandleFunc("/logs/detail", h.LogDetailPage)

	// API
	mux.HandleFunc("/api/logs", h.LogsAPI)
	mux.HandleFunc("/requests", h.RequestsPage)
	mux.HandleFunc("/api/requests", h.RequestsAPI)

	server := &http.Server{
		Addr:              addr,
		Handler:           requestLogger(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("loki-ui listening on %s", addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		log.Printf("loki-ui shutdown requested")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		log.Printf("loki-ui shutdown completed")
		return nil
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if r.Method == http.MethodHead {
		return
	}

	_, _ = w.Write([]byte("ok\n"))
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		requestID := newRequestID()

		w.Header().Set("X-Request-Id", requestID)

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		writeRequestLog(requestLogEntry{
			Type:       "loki_ui_request",
			RequestID:  requestID,
			Method:     r.Method,
			Path:       r.URL.Path,
			Status:     rec.status,
			DurationMS: time.Since(started).Milliseconds(),
		})
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

type requestLogEntry struct {
	Type       string `json:"type"`
	RequestID  string `json:"request_id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	DurationMS int64  `json:"duration_ms"`
}

func writeRequestLog(entry requestLogEntry) {
	b, err := json.Marshal(entry)
	if err != nil {
		log.Printf(`{"type":"loki_ui_request","error":"encode_failed"}`)
		return
	}

	log.Print(string(b))
}

func newRequestID() string {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "request-id-unavailable"
	}

	return hex.EncodeToString(b[:])
}
