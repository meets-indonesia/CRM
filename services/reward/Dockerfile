FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o reward main.go

FROM alpine:latest
WORKDIR /app
RUN mkdir -p /app/uploads
COPY --from=builder /app/reward .
EXPOSE 8084
CMD ["./reward"]