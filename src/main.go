package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	app := &application{
		templates: tmpl,
		state:     &appState{},
	}

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle root path
	http.HandleFunc("/", app.counterHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
