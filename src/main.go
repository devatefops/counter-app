package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"sync"

	// Import the content package to access embedded files
	"counter-app/src/content"
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

// counterHandler handles HTTP requests.
func (app *application) counterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	app.state.mu.Lock()
	defer app.state.mu.Unlock()

	// Handle form submissions
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
	// Load templates from the correct path relative to src/
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	// Parse templates from the embedded filesystem for robustness.
	templates, err := template.ParseFS(content.TemplatesFS, "templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	app := &application{
		templates: tmpl,
		templates: templates,
		state:     &appState{},
	}

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Serve static files from the embedded filesystem.
	staticFS, err := fs.Sub(content.StaticFS, "static")
	if err != nil {
		log.Fatalf("failed to create static sub-filesystem: %v", err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Handle root path
	http.HandleFunc("/", app.counterHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
