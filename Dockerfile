# builder stage
FROM golang:1.24 AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]