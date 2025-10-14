FROM node:20-alpine AS builder
WORKDIR /app

COPY web/package*.json ./
RUN npm install

COPY web .
RUN npm run build

FROM node:20-alpine AS production
WORKDIR /app

LABEL org.opencontainers.image.title="Uptimatic"
LABEL org.opencontainers.image.description="Uptimatic frontend service (Next.js)"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.authors="Nandan <nandanramdani608@gmail.com>"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.url="https://uptimatic.aeria.my.id"
LABEL org.opencontainers.image.source="https://github.com/naandan/uptimatic"
LABEL org.opencontainers.image.documentation="https://docs.uptimatic.aeria.my.id"

ENV NODE_ENV=production

RUN addgroup -S uptimatic && adduser -S -G uptimatic uptimatic

COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

RUN chown -R uptimatic:uptimatic /app

USER uptimatic

EXPOSE 3000

CMD ["npm", "start"]