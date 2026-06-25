FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json ./
RUN npm install --omit=dev

FROM node:20-alpine
WORKDIR /app
RUN addgroup -S collab && adduser -S collab -G collab
COPY --from=deps /app/node_modules ./node_modules
COPY server.js ./
USER collab
EXPOSE 1234
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:1234/health 2>/dev/null || exit 0
CMD ["node", "server.js"]
