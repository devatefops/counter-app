# Counter Platform

A simple event-driven system built with Go, consisting of two services:

* **Counter Service** → maintains a counter and exposes an API
* **Notifier Service** → consumes the counter value and reacts to changes

This project demonstrates basic microservices design, testing, and CI/CD readiness.

---

## 🏗️ Architecture

```id="7b8j3p"
+-------------------+        HTTP        +--------------------+
|  Counter Service  |  <-------------->  |  Notifier Service  |
|                   |                    |                    |
| - Web UI          |                    | - Polls API        |
| - REST API        |                    | - Reacts to value  |
| - In-memory state |                    |                    |
+-------------------+                    +--------------------+
```

---

## 📦 Project Structure

```id="v8qv9c"
.
├── counter-service/
│   ├── cmd/
│   ├── templates/
│   └── go.mod
│
├── notifier-service/
│   ├── cmd/
│   └── go.mod
│
└── README.md
```

---

## 🚀 Services

### 1. Counter Service

* Web UI to increment/reset counter
* Thread-safe using mutex
* Exposes API:

```id="hcbk0q"
GET /api/counter
```

Response:

```json id="v5gh9m"
{ "value": 3 }
```

---

### 2. Notifier Service

* Periodically calls counter API
* Detects changes in value
* Can be extended to:

  * send logs
  * trigger alerts
  * integrate with external systems

---

## ▶️ Run Locally

### Start Counter Service

```bash id="6bsk9w"
cd counter-service
go run ./cmd
```

Runs on:

```id="rqv1ps"
http://localhost:8080
```

---

### Start Notifier Service

```bash id="p9x4tb"
cd notifier-service
go run ./cmd
```

---

## 🧪 Run Tests

From each service:

```bash id="6f0h5z"
go test -race -v ./...
```

---

## 🐳 Docker (optional)

Each service can be containerized independently:

```bash id="j4j9r8"
docker build -t counter-service ./counter-service
docker build -t notifier-service ./notifier-service
```

---

## ⚙️ CI/CD

Each service is tested independently using GitHub Actions:

* Runs unit tests
* Ensures code quality before build
* Can be extended to:

  * build Docker images
  * push to registry
  * deploy to Kubernetes

---

## 🧠 Design Principles

* Separation of concerns (each service has a single responsibility)
* Stateless communication via HTTP
* Testable architecture (no filesystem dependency in tests)
* Minimal dependencies (standard library)

---

## 🔮 Future Improvements

* Add message broker (Kafka / RabbitMQ)
* Replace polling with event-driven communication
* Add observability (Prometheus + Grafana)
* Deploy using Kubernetes
* Add distributed tracing

---

## 📄 License

MIT
