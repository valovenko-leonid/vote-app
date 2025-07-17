#!/bin/sh
set -e

CERT_PATH="/etc/letsencrypt/live/$DOMAIN/fullchain.pem"

# Сначала запустим nginx с pre-cert конфигом
echo ">>> Запускаем nginx с pre-cert конфигом..."
cp /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.bak
cp /etc/nginx/conf.d/pre-cert.conf /etc/nginx/conf.d/default.conf
nginx &

# Ждём, пока nginx поднимется
sleep 5

if [ "$ENV" = "production" ] && [ ! -f "$CERT_PATH" ]; then
  echo ">>> Получаем сертификат..."
  certbot certonly \
    --webroot -w /var/www/certbot \
    --email "$EMAIL" --agree-tos --no-eff-email \
    -d "$DOMAIN"
fi

# Останавливаем nginx и запускаем с боевым конфигом
echo ">>> Перезапускаем nginx с production конфигом..."
nginx -s stop
cp /etc/nginx/conf.d/production.conf /etc/nginx/conf.d/default.conf
nginx -g "daemon off;"