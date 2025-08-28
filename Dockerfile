# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bailanysta-api ./api/cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bailanysta-api .

# Copy migrations
COPY --from=builder /app/api/internal/db/migrations ./migrations

# Expose port
EXPOSE 8080

# Set timezone
ENV TZ=UTC

# Run the binary
CMD ["./bailanysta-api"]
