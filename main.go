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

	errCh := make(chan error)
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Start server on port %s", port)
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		log.Fatalf("Server error, terminating: %+v", err)
	case sig := <-sigCh:
		log.Printf("Normal shutdown by signal: %+v", sig)
	}

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Stopping server: %+v", err)
	}

	os.Exit(0)
}
