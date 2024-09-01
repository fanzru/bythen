# Start with the official Go 1.22.5 image as a base image
FROM golang:1.22.5-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apk update && apk add --no-cache git

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code to the workspace
COPY . .

# Build the Go application
RUN go build -o main cmd/main.go

# Use a minimal image for the final build
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

# # Copy .env file if any
# COPY .env ./

# Copy the Swagger documentation files
COPY docs/swagger /app/docs/swagger

# Expose the application port
EXPOSE 8080

# Run the Go application
CMD ["./main"]
