package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// appState holds the counter and its mutex.
type appState struct {
	mu      sync.Mutex
	counter int
}

// application holds the application's dependencies.
type application struct {
	templates *template.Template
	state     *appState
}

func (app *application) counterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	app.state.mu.Lock()
	defer app.state.mu.Unlock()

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		switch r.FormValue("action") {
		case "increment":
			app.state.counter++
		case "reset":
			app.state.counter = 0
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := app.templates.ExecuteTemplate(w, "index.html", app.state.counter); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	// Get working directory (project root if run correctly)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templatesPath := filepath.Join(wd, "templates", "index.html")
	staticPath := filepath.Join(wd, "static")

	app := &application{
		templates: template.Must(template.ParseFiles(templatesPath)),
		state:     &appState{},
	}

	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.counterHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
