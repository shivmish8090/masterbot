FROM golang:1.24-bullseye AS builder

WORKDIR /app

COPY go.* ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app .

# Final minimal image
FROM scratch

WORKDIR /app

# Optional: Add CA certs for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/app .

ENTRYPOINT ["/app/app"]