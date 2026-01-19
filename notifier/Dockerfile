# -------- Build Stage --------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install Git (needed for Go modules)
RUN apk add --no-cache git

# Copy the entire cmd folder (contains main.go, go.mod, go.sum)
COPY cmd ./cmd

# Switch to the cmd directory
WORKDIR /app/cmd

# Download module dependencies
RUN go mod download

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/notifier

# -------- Runtime Stage --------
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/notifier .

# Run the binary
CMD ["./notifier"]
