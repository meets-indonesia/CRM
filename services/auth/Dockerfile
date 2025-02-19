# Stage 1: Build
FROM golang:alpine AS builder
WORKDIR /app
# Install dependencies for building
RUN apk add --no-cache gcc musl-dev
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy source code
COPY . .
# Build the application statically
RUN go build -o main .

# Stage 2: Minimal runtime
FROM alpine:latest
WORKDIR /root/
# Install necessary runtime dependencies (if any)
RUN apk --no-cache add ca-certificates
# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .
# Copy templates folder from builder stage
COPY --from=builder /app/templates ./templates
# Expose the application port
EXPOSE 8080
# Run the application
CMD ["./main"]