# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Set working directory
# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for caching not used for now 
#COPY go.mod go.sum ./
#RUN go mod download
# Copy go.mod and go.sum files to download dependencies.
# This leverages Docker's layer caching for faster builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
# Build the Go app as a static binary. This is crucial for the scratch image.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./src/main.go

# Stage 2: Create the final image (using lightweight Alpine)
FROM alpine:3.18
# Stage 2: Create the final, minimal image from scratch
FROM scratch

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose port
# Expose the port the app runs on
EXPOSE 8080

# Run the app
# Command to run the executable
CMD ["./main"]
