FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o feedback main.go

FROM alpine:latest
WORKDIR /app
# Create upload directory
RUN mkdir -p /app/uploads
COPY --from=builder /app/feedback .
EXPOSE 8083
CMD ["./feedback"]