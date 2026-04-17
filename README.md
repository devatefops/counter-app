# Counter Service

A simple Go web application that maintains a counter and exposes both a UI and a JSON API.

---

## 🚀 Features

* Increment and reset counter via web UI
* Thread-safe state using mutex
* JSON API endpoint for integrations
* Designed for CI/CD and containerization

---

## 🏗️ Project Structure

```
.
├── cmd/
│   ├── main.go
│   └── main_test.go
├── templates/
│   └── index.html
├── static/
├── go.mod
```

---

## ▶️ Run Locally

### 1. Install Go (>= 1.22)

### 2. Run the app

```bash
go run ./cmd
```

App will start on:

```
http://localhost:8080
```

---

## 🧪 Run Tests

```bash
go test -race -v ./...
```

---

## 🌐 API Endpoints

### GET `/`

* Returns HTML page with counter

### POST `/`

* Form actions:

  * `increment`
  * `reset`

### GET `/api/counter`

```json
{
  "value": 3
}
```

---

## 🐳 Run with Docker

### Build image

```bash
docker build -t counter-service .
```

### Run container

```bash
docker run -p 8080:8080 counter-service
```

---

## ⚙️ CI/CD (GitHub Actions)

Example workflow:

```yaml
- uses: actions/checkout@v4

- uses: actions/setup-go@v5
  with:
    go-version: '1.22'

- run: go test -race ./...
```

---

## 🧠 Design Notes

* Uses `net/http` (no external frameworks)
* Templates are injected → easy to test
* Unit tests avoid filesystem dependency
* Safe for concurrent access

---

## 🔮 Future Improvements

* Add persistent storage (Redis / DB)
* Add metrics (Prometheus)
* Add authentication
* Kubernetes deployment

---

## 📄 License

MIT
