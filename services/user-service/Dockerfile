# --- Build Stage ---
# Use the official Go image as the builder.
# Pinning to a specific version ensures consistent builds.
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
# Download dependencies. This is cached if the mod/sum files don't change.
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app.
# The -o flag specifies the output file name.
# CGO_ENABLED=0 is important for creating a static binary that can run in a minimal image.
# -ldflags="-w -s" strips debug information, reducing the binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /user-service .

# --- Final Stage ---
# Use a minimal, non-root base image for the final container.
# Alpine is small and has a package manager if needed for debugging.
FROM alpine:3.18

# Set the working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /user-service /user-service

# Expose the port the app runs on (as defined in main.go)
EXPOSE 3000

# This is the command that will be executed when the container starts
ENTRYPOINT ["/user-service"]
