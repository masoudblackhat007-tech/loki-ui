package httpserver

import (
	"context"
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

	// UI
	mux.HandleFunc("/logs", h.LogsPage)
	mux.HandleFunc("/logs/detail", h.LogDetailPage)

	// API
	mux.HandleFunc("/api/logs", h.LogsAPI)
	mux.HandleFunc("/requests", h.RequestsPage)
	mux.HandleFunc("/api/requests", h.RequestsAPI)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
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
