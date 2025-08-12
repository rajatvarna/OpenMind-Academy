# --- Build Stage ---
# Use the official Go image as the builder.
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app, creating a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /content-service .

# --- Final Stage ---
# Use a minimal base image for the final container.
FROM alpine:3.18

# Set the working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /content-service /content-service

# Expose the port the app runs on (as defined in main.go)
EXPOSE 3001

# This is the command that will be executed when the container starts
ENTRYPOINT ["/content-service"]
