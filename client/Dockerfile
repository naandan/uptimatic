FROM node:20-alpine AS builder
WORKDIR /app

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY web/pnpm-lock.yaml web/package.json ./
RUN pnpm install --frozen-lockfile

COPY web ./
RUN pnpm run build

FROM node:20-alpine AS production
WORKDIR /app

LABEL org.opencontainers.image.title="Uptimatic" \
      org.opencontainers.image.description="Uptimatic frontend service (Next.js)" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.authors="Nandan <nandanramdani608@gmail.com>" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.url="https://uptimatic.aeria.my.id" \
      org.opencontainers.image.source="https://github.com/naandan/uptimatic" \
      org.opencontainers.image.documentation="https://docs.uptimatic.aeria.my.id"

ENV NODE_ENV=production

RUN corepack enable && corepack prepare pnpm@latest --activate
RUN addgroup -S uptimatic && adduser -S -G uptimatic uptimatic

COPY web/pnpm-lock.yaml web/package.json ./
RUN pnpm install --prod --frozen-lockfile

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

RUN chown -R uptimatic:uptimatic /app
USER uptimatic

EXPOSE 3000

CMD ["pnpm", "start"]