# ---- Build stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o hse-server main.go

# ---- Run stage ----
FROM alpine:3.20
WORKDIR /app


RUN apk add --no-cache libreoffice-writer ttf-dejavu fontconfig

COPY --from=builder /app/hse-server .
COPY --from=builder /app/templates ./templates

EXPOSE 3000
CMD ["./hse-server"]