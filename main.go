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

func (c CGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	scriptPath := c.cleanPath(r.URL.Path)

	if strings.HasSuffix(scriptPath, "/") {
		c.statusError(w, http.StatusForbidden)
		return
	}

	scriptFile := c.ScriptsDir + scriptPath

	stat, err := os.Stat(scriptFile)
	if err != nil || stat.IsDir() {
		c.statusError(w, http.StatusNotFound)
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

func (c CGIHandler) statusError(w http.ResponseWriter, statusCode int) {
	statusText := http.StatusText(statusCode)
	w.WriteHeader(statusCode)
	w.Write([]byte(statusText))
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
