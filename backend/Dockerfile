# Stage 1: Build the Go binary
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go files to the container
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Build the Go app
RUN go build -o wasaBe

# Stage 2: Run the binary using a lightweight image
FROM debian:bullseye-slim

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/wasaBe .

# Expose the port your app listens on
EXPOSE 8080

# Run the binary
CMD ["./wasaBe"]
