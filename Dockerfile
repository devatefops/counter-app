# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
# This leverages Docker's layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app. CGO_ENABLED=0 creates a statically linked binary.
# This is crucial for running in a minimal 'scratch' image.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./src/main.go

# Stage 2: Create the final, minimal image
FROM scratch

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

