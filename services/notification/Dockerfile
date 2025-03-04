FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o notification main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notification .
EXPOSE 8087
CMD ["./notification"]