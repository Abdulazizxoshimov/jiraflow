#!/usr/bin/env bash
# First-time Let's Encrypt certificate setup.
# Run ONCE on a fresh server before starting the full stack.
#
# Usage:
#   DOMAIN=yourdomain.com EMAIL=admin@yourdomain.com ./scripts/init-letsencrypt.sh

set -euo pipefail

DOMAIN="${DOMAIN:?Set DOMAIN=yourdomain.com}"
EMAIL="${EMAIL:?Set EMAIL=admin@yourdomain.com}"
COMPOSE="docker compose -f deployments/docker-compose.yml"

echo "[init] Starting HTTP-only nginx for ACME challenge..."
$COMPOSE up -d nginx

echo "[init] Requesting certificate for ${DOMAIN}..."
$COMPOSE run --rm certbot certonly \
  --webroot \
  --webroot-path /var/www/certbot \
  --email "$EMAIL" \
  --agree-tos \
  --no-eff-email \
  -d "$DOMAIN"

echo "[init] Reloading nginx with HTTPS..."
$COMPOSE exec nginx nginx -s reload

echo "[init] Done. Certificate stored in certbot_certs volume."
echo "[init] Auto-renewal runs every 12h via the certbot service."
