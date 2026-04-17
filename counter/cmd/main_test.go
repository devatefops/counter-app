package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setupApp() *application {
	tmpl := template.Must(template.New("test").Parse(`
		<html><body>{{.}}</body></html>
	`))

	return &application{
		templates: tmpl,
		state:     &appState{},
	}
}

func TestCounterHandler_Get(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.counterHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}
func TestCounterHandler_Increment(t *testing.T) {
	app := setupApp()

	form := url.Values{}
	form.Add("action", "increment")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.counterHandler)

	handler.ServeHTTP(rr, req)

	if app.state.counter != 1 {
		t.Errorf("expected counter 1, got %d", app.state.counter)
	}
}
func TestCounterHandler_Reset(t *testing.T) {
	app := setupApp()
	app.state.counter = 5

	form := url.Values{}
	form.Add("action", "reset")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.counterHandler)

	handler.ServeHTTP(rr, req)

	if app.state.counter != 0 {
		t.Errorf("expected counter 0, got %d", app.state.counter)
	}
}
func TestAPICounterHandler(t *testing.T) {
	app := setupApp()
	app.state.counter = 3

	req := httptest.NewRequest(http.MethodGet, "/api/counter", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.apiCounterHandler)
	handler.ServeHTTP(rr, req)

	expected := `{"value": 3}`

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected %s, got %s", expected, rr.Body.String())
	}
}
