#!/bin/sh
set -e

CERT_PATH="/etc/letsencrypt/live/$DOMAIN/fullchain.pem"
PRECERT_CONF="/etc/nginx/pre-cert.conf"
DEV_CONF="/etc/nginx/development.conf"
PROD_CONF="/etc/nginx/production.conf"

rm -f /etc/nginx/conf.d/*.conf
cp "$PRECERT_CONF" /etc/nginx/conf.d/default.conf

# Не production → запускаем дев-конфиг
if [ "$ENV" != "production" ]; then
  echo ">>> dev-окружение. Запускаем nginx..."
  cp "$DEV_CONF" /etc/nginx/conf.d/default.conf

  echo ">>> Список .conf-файлов в /etc/nginx/conf.d:"
ls -la /etc/nginx/conf.d/
cat /etc/nginx/conf.d/default.conf
  nginx -g "daemon off;"
  exit 0
fi

# Production, проверка сертификата
if [ -f "$CERT_PATH" ]; then
  echo ">>> Сертификат найден. Запускаем nginx с production.conf..."
  cp "$PROD_CONF" /etc/nginx/conf.d/default.conf
else
  echo ">>> Сертификат НЕ найден. Запускаем nginx с pre-cert.conf..."
  cp "$PRECERT_CONF" /etc/nginx/conf.d/default.conf
fi


echo ">>> Список .conf-файлов в /etc/nginx/conf.d:"
ls -la /etc/nginx/conf.d/
cat /etc/nginx/conf.d/default.conf
nginx -g "daemon off;"