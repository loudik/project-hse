# ---- Build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o hse-server main.go

# ---- Run stage ----
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/hse-server .

EXPOSE 3000
CMD ["./hse-server"]
