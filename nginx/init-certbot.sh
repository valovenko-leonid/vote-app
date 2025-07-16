#!/bin/sh

set -e

CERT_PATH="/etc/letsencrypt/live/$DOMAIN/fullchain.pem"

if [ "$ENV" = "production" ] && [ ! -f "$CERT_PATH" ]; then
  echo ">>> Сертификат не найден. Запускаем Certbot..."
  
  certbot certonly \
    --webroot -w /var/www/certbot \
    --email "$EMAIL" \
    --agree-tos --no-eff-email \
    -d "$DOMAIN"

  echo ">>> Сертификат получен. Перезапускаем nginx..."
  nginx -s reload
else
  echo ">>> Сертификат уже существует или не production. Запускаем nginx..."
fi

# Запускаем nginx в foreground
nginx -g "daemon off;"