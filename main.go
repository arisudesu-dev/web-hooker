package main

import (
	"context"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

type CGIHandler struct {
	ScriptsDir string
}

func (c CGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REQ] %+v", r.URL)

	scriptPath := c.cleanPath(r.URL.Path)

	if strings.HasSuffix(scriptPath, "/") {
		if err := c.statusError(w, http.StatusForbidden); err != nil {
			log.Printf("[ERR] %+v", err)
		}
		return
	}

	scriptFile := filepath.FromSlash(c.ScriptsDir + scriptPath)
	log.Printf("[INF] %+v => %+v", r.URL, scriptFile)

	stat, err := os.Stat(scriptFile)
	if err != nil || stat.IsDir() {
		if err := c.statusError(w, http.StatusNotFound); err != nil {
			log.Printf("[ERR] %+v", err)
		}
		return
	}

	cgiHandler := cgi.Handler{
		Root: scriptPath,
		Path: scriptFile,
		Dir:  c.ScriptsDir,
	}
	cgiHandler.ServeHTTP(w, r)
}

func (c CGIHandler) cleanPath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	np := path.Clean(p)

	if strings.HasSuffix(p, "/") && np != "/" {
		np += "/"
	}
	return np
}

func (c CGIHandler) statusError(w http.ResponseWriter, statusCode int) error {
	statusText := http.StatusText(statusCode)
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(statusText))
	return err
}

func main() {
	scriptsDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot obtain current working directory: %+v", err)
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
