# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying the source code
COPY . .

# Building a binary
RUN go build -o url-shortener ./cmd/main.go

# ---------- Run stage ----------
FROM alpine:latest

WORKDIR /app

# Copying the binary
COPY --from=builder /app/url-shortener .

# Exposing the port
EXPOSE 8080

# Launching the server
CMD ["./url-shortener"]