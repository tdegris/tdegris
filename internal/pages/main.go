package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	port     = flag.Int("port", 8080, "http port")
	local    = flag.Bool("local", true, "local connections only")
	logQuery = flag.Bool("logq", true, "log queries")
)

func runGoGenerate() error {
	wdir, err := os.Getwd()
	if err != nil {
		return err
	}
	moduleRoot, err := findParent("go.mod")
	if err != nil {
		return err
	}
	if err := os.Chdir(moduleRoot); err != nil {
		return err
	}
	defer func() {
		os.Chdir(wdir)
	}()
	cmd := exec.Command("go", "generate", "./...")
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Println(string(out))
	}
	if err != nil {
		return err
	}
	return nil
}

func buildAddr() string {
	addr := fmt.Sprintf(":%d", *port)
	if *local {
		addr = "localhost" + addr
	}
	return addr
}

func findParent(folder string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for cwd != "/" {
		modPath := filepath.Join(cwd, folder)
		if _, err := os.Stat(modPath); err == nil {
			return cwd, nil
		}
		cwd = filepath.Dir(cwd)
	}
	return "", fmt.Errorf("parent %s not found", folder)
}

func mainHandler(fs http.FileSystem, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/res/main.wasm" {
		if err := runGoGenerate(); err != nil {
			http.Error(w, fmt.Sprintf("cannot generate WASM file: %v", err), http.StatusInternalServerError)
			return
		}
	}
	http.FileServer(fs).ServeHTTP(w, r)
}

func run() error {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	if *logQuery {
		r.Use(middleware.Logger)
	}
	projectRoot, err := findParent(".git")
	if err != nil {
		return err
	}
	projectFS := http.FS(os.DirFS(projectRoot))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		mainHandler(projectFS, w, r)
	})
	addr := buildAddr()
	fmt.Printf("Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return fmt.Errorf("cannot run HTTP server: %v", err)
	}
	return nil
}

func main() {
	flag.Parse()
	err := run()
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
}
