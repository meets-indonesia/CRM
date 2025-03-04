FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o article main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/article .
EXPOSE 8086
CMD ["./article"]