FROM golang:1.24.6-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o uptimatic .

FROM alpine:3.18
WORKDIR /app

LABEL org.opencontainers.image.title="Uptimatic" \
      org.opencontainers.image.description="Uptimatic backend service (Go + Supervisor)" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.authors="Nandan <nandanramdani608@gmail.com>" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.url="https://uptimatic.aeria.my.id" \
      org.opencontainers.image.source="https://github.com/naandan/uptimatic" \
      org.opencontainers.image.documentation="https://docs.uptimatic.aeria.my.id"

RUN apk add --no-cache supervisor curl tzdata

RUN addgroup -S uptimatic && adduser -S -G uptimatic uptimatic

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz \
    | tar xz -C /usr/local/bin migrate

COPY --from=builder /app/uptimatic /app/uptimatic
COPY supervisord.conf /etc/supervisord.conf
COPY supervisor /etc/supervisor/conf.d
COPY app/migrations /app/migrations
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh /app/uptimatic && \
    mkdir -p /var/log && \
    chown -R uptimatic:uptimatic /var/log /app && \
    chmod -R 755 /var/log /app

USER uptimatic

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["supervisord", "-n", "-c", "/etc/supervisord.conf"]