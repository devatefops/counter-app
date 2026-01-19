# -------- Build Stage --------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Git for module dependencies
RUN apk add --no-cache git

# Copy go.mod & go.sum for caching
# COPY go.mod go.sum ./
# RUN go mod download

# Copy all source code
COPY . .

# Build the Go binary from cmd package
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notifier ./cmd

# -------- Runtime Stage --------
FROM alpine:3.19

WORKDIR /app

# Certificates (for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/notifier .

# Run the binary
CMD ["./notifier"]
