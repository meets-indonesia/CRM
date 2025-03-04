FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o inventory main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/inventory .
EXPOSE 8085
CMD ["./inventory"]