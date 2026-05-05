package main

import (
	"log"
	"net"
	"os"
	"strings"

	"loki-ui/internal/httpserver"
)

func main() {
	addr := strings.TrimSpace(os.Getenv("LISTEN_ADDR"))
	if addr == "" {
		addr = "127.0.0.1:18090"
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil || port == "" {
		log.Fatalf(
			"invalid LISTEN_ADDR=%q (expected 127.0.0.1:port or localhost:port)",
			addr,
		)
	}

	// --- reject wildcard binds (:18090, 0.0.0.0, ::) ---
	if host == "" || host == "0.0.0.0" || host == "::" {
		log.Fatalf(
			"refusing to listen on wildcard address: %q (must be loopback only)",
			addr,
		)
	}

	// --- allow localhost explicitly ---
	if host == "localhost" {
		start(addr)
		return
	}

	// --- strict IP validation ---
	ip := net.ParseIP(host)
	if ip == nil {
		log.Fatalf(
			"invalid LISTEN_ADDR host=%q (must be 127.0.0.1 or localhost)",
			host,
		)
	}

	// must be loopback AND exactly 127.0.0.1
	if !ip.IsLoopback() || host != "127.0.0.1" {
		log.Fatalf(
			"refusing to listen on non-loopback address: %q (only 127.0.0.1 or localhost allowed)",
			addr,
		)
	}

	start(addr)
}

func start(addr string) {
	log.Printf("loki-ui listening on %s", addr)

	if err := httpserver.Start(addr); err != nil {
		log.Fatal(err)
	}
}
