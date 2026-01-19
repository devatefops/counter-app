# -------- Build Stage --------
FROM golang:1.22-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary from the cmd package
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/notifier ./cmd

# -------- Runtime Stage --------
FROM alpine:3.19

WORKDIR /app

# Add CA certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/notifier .

# Run the binary
CMD ["./notifier"]
