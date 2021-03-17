package main

import (
	"context"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

type CGIHandler struct {
	ScriptsDir string
}

func (c CGIHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// strip leading '/' and combine with scripts path
	scriptPath := path.Join(c.ScriptsDir, req.URL.Path[1:])

	// check if:
	//   scriptPath is inside scriptsDir;
	//   scriptPath is not empty;
	//   scriptPath is a file (doesn't end with '/')
	if !strings.HasPrefix(scriptPath, c.ScriptsDir) ||
		strings.TrimPrefix(scriptPath, c.ScriptsDir) == "" ||
		strings.HasSuffix(scriptPath, "/") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	cgiHandler := cgi.Handler{
		Path: scriptPath,
		Dir:  c.ScriptsDir,
	}
	cgiHandler.ServeHTTP(rw, req)
}

func main() {
	scriptsDir, ok := os.LookupEnv("SCRIPTS_DIR")
	if !ok {
		scriptsDir = "/"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	server := http.Server{}
	server.Addr = ":" + port
	server.Handler = &CGIHandler{ScriptsDir: scriptsDir}

	brkCh := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Signal received: %+v", <-sigCh)
		_ = server.Shutdown(context.Background())
		close(brkCh)
	}()

	log.Printf("Start server on port %s", port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error, terminating: %+v", err)
	}

	log.Printf("Waiting for server shutdown...")
	<-brkCh
	os.Exit(0)
}
