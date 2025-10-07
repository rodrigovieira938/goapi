# ================================
# Stage 1: Build the Go binary
# ================================
FROM golang:1.25-alpine AS builder

# Install build tools
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum separately (to leverage Docker cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate cmd/migrate/main.go

# ================================
# Stage 2: Minimal runtime image
# ================================
FROM alpine:3.20

WORKDIR /app

# Copy only the binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .

# Copy needed files
COPY ./.env .
COPY ./migrations ./migrations

# Expose the port your Go app listens on
EXPOSE 8080

# Start the app
CMD ["sh", "-c", "./migrate up && ./server"]
