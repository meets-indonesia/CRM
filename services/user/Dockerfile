FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o user main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user .
EXPOSE 8082
CMD ["./user"]