# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum* ./
RUN go mod download

# Copy application source
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/webservice ./cmd/server/

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/webservice .

# Copy the .env file
COPY .env .

# Expose port (this is just documentation, you'll need to define the port in .env)
EXPOSE 8080

# Run the web service
CMD ["./webservice"]