# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app

# Install dependencies needed for building
RUN apk add --no-cache git gcc musl-dev

# Copy Go modules manifest and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8083

# Run the application
CMD ["./main"]