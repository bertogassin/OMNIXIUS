# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api .

# Run stage (minimal; SQLite needs libc)
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /api .
# Default DB path in container; override with -e DB_PATH=/data/omnixius.db
ENV PORT=8080
EXPOSE 8080
VOLUME /data
CMD ["./api"]
