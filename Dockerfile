# Use the official Golang image as a base image
FROM golang:1.17 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN go build -o main cmd/main.go

# Start a new stage from the official Golang image
FROM golang:1.17

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main

# Expose the API port
EXPOSE 8080

# Run the binary
CMD ["/app/main"]
