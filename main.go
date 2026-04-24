package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mostlygeek/llama-swap/proxy"
)

const version = "0.0.1"

func main() {
	var (
		configPath  string
		listenAddr  string
		showVersion bool
		logRequests bool
	)

	flag.StringVar(&configPath, "config", "config.yaml", "path to configuration file")
	// Changed default port from 8080 to 11434 to match Ollama's default port,
	// making it easier to use as a drop-in replacement in my local setup.
	flag.StringVar(&listenAddr, "listen", ":11434", "address to listen on (host:port)")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	// Enabling request logging by default so I can see what's hitting the proxy
	// during development without needing to remember to pass the flag each time.
	flag.BoolVar(&logRequests, "log-requests", true, "log all incoming requests")
	flag.Parse()

	if showVersion {
		fmt.Printf("llama-swap version %s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := proxy.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config from %q: %v", configPath, err)
	}

	log.Printf("llama-swap v%s starting", version)
	log.Printf("listening on %s", listenAddr)
	log.Printf("loaded %d model(s) from config", len(cfg.Models))
	// Print a reminder about the config file location so it's easy to spot in logs.
	log.Printf("using config file: %s", configPath)

	// Create and start the proxy server
	server, err := proxy.NewServer(cfg, proxy.ServerOptions{
		ListenAddr:  listenAddr,
		LogRequests: logRequests,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
