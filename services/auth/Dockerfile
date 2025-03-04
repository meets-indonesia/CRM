FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o auth main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth .
EXPOSE 8081
CMD ["./auth"]