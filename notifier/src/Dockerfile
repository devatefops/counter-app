# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notifier ./src

# Final Stage
FROM alpine:latest

WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/notifier .
CMD ["./notifier"]
