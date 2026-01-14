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

// counterHandler handles the HTTP requests to our web app.
func (app *application) counterHandler(w http.ResponseWriter, r *http.Request) {
	// We only want to handle requests for the root path.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	app.state.mu.Lock()

	// Handle form submission for buttons
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			// Log the error and send a bad request response.
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			app.state.mu.Unlock()
			return
		}
		switch r.FormValue("action") {
		case "increment":
			app.state.counter++
		case "reset":
			app.state.counter = 0
		}
	}

	// It's good practice to release the lock before doing I/O (writing the response)
	currentCount := app.state.counter
	app.state.mu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := app.templates.ExecuteTemplate(w, "index.html", currentCount)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Initialize the application dependencies.
	app := &application{
		templates: template.Must(template.ParseFiles("../templates/index.html")),
		state:     &appState{},
	}

	// Serve static files (like CSS)
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.counterHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
