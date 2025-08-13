# --- Build Stage ---
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /forum-service .

# --- Final Stage ---
FROM alpine:3.18
WORKDIR /
COPY --from=builder /forum-service /forum-service
EXPOSE 3006
ENTRYPOINT ["/forum-service"]
