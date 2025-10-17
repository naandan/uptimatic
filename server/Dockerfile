FROM golang:1.24.6-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY app .
RUN go mod download
RUN go build -o uptimatic .

FROM alpine:3.18
WORKDIR /app

LABEL org.opencontainers.image.title="Uptimatic"
LABEL org.opencontainers.image.description="Uptimatic backend service (Go + Supervisor)"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.authors="Nandan <nandanramdani608@gmail.com>"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.url="https://uptimatic.aeria.my.id"
LABEL org.opencontainers.image.source="https://github.com/naandan/uptimatic"
LABEL org.opencontainers.image.documentation="https://docs.uptimatic.aeria.my.id"

RUN apk add --no-cache supervisor bash curl shadow tzdata

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
