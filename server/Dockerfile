# Stage 1: Build Go binary
FROM golang:1.24.6-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o uptimatic .

# Stage 2: Runtime
FROM alpine:3.18
WORKDIR /app

# Install s6-overlay
ENV S6_OVERLAY_VERSION=v3.1.2.1
ADD https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64.tar.gz /tmp/
RUN tar xzf /tmp/s6-overlay-amd64.tar.gz -C / && rm /tmp/s6-overlay-amd64.tar.gz

# Copy binary
COPY --from=builder /app/uptimatic /app/uptimatic

# Copy s6 service folder
COPY s6 /etc/services.d
COPY entrypoint.sh /

# Create log folders
RUN mkdir -p /var/log/server /var/log/worker /var/log/scheduler \
    && chmod -R 755 /var/log \
    && chmod +x /app/uptimatic \
    && chmod +x /entrypoint.sh \
    && chmod -R +x /etc/services.d

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
