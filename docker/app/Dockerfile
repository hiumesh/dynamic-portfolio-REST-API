# Production-ready Dockerfile for Go REST API with Gin framework

# Step 1: Build stage
FROM golang:latest AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Create and set the working directory
WORKDIR /app

# Cache module dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o server .

# Step 2: Final stage
FROM gcr.io/distroless/static:nonroot

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Expose the port the Gin app will run on
EXPOSE 8080

# Use a non-root user for better security
USER nonroot:nonroot

# Command to run the application
CMD ["/app/server"]
