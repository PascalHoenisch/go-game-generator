# Use a minimal base image for Go
FROM golang:1.20-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum first for dependency resolution
COPY app/go.mod ./
# If there’s a go.sum, uncomment and copy it too to ensure reproducible builds
# COPY go.sum ./
RUN go mod download

# Copy the entire project (including main.go and subdirectories)
COPY app .

# Build the Go application
RUN go build -o app main.go

# Create a smaller runtime image
FROM alpine:latest

# Set working directory inside the runtime container
WORKDIR /app

# Copy the compiled Go application from the builder image
COPY --from=builder /app/app .

# Copy other necessary files (e.g., templates folder)
COPY --from=builder /app/templates ./templates

# Expose port 80 for the Go application
EXPOSE 80

# Command to run the compiled Go application
CMD ["./app"]