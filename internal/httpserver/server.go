package httpserver

import (
	"log"
	"net/http"
	"time"
)

func Start(addr string) error {
	h := NewHandler()

	mux := http.NewServeMux()

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

	log.Printf("loki-ui listening on %s", addr)
	return server.ListenAndServe()
}
