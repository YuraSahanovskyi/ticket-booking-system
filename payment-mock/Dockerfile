FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/payment-server .
COPY templates/ ./templates/

EXPOSE 8081
CMD ["./payment-server"]