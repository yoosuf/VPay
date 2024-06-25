# Use the official Go image as a build stage
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main cmd/main.go

# Use a minimal image for the runtime stage
FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache bash

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled Go binary from the build stage
COPY --from=builder /app/main /app/main

# Copy the TLS certificates and wait-for-it.sh script
COPY server.crt server.key /app/
COPY wait-for-it.sh /app/

# Make the wait-for-it.sh script executable
RUN chmod +x /app/wait-for-it.sh

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["/app/wait-for-it.sh", "db", "3306", "--", "/app/main"]
