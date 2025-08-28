# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bailanysta-api ./api/cmd/api/

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bailanysta-api .

# Copy migrations
COPY --from=builder /app/api/internal/db/migrations ./api/internal/db/migrations

# Expose port
EXPOSE 8080

# Command to run
CMD ["./bailanysta-api"]
