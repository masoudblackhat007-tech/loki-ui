package httpserver

import (
	"log"
	"net/http"
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

	log.Printf("loki-ui listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}
