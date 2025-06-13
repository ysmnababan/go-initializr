# 🧱 Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Set working directory in the builder image
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies (cached if mod files don't change)
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN go build -o main .

# 🪶 Stage 2: runtime — also needs Go SDK
FROM golang:1.24-alpine

# Set working directory in the final image
WORKDIR /app

# Copy only the built binary from the builder image
COPY --from=builder /app/main .

# Copy the template files
COPY template/ ./template/

# Document the port (for clarity and tooling)
EXPOSE 1323

# Run the binary when the container starts
CMD ["./main"]
