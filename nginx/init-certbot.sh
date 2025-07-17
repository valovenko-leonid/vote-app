#!/bin/sh
set -e

CERT_PATH="/etc/letsencrypt/live/$DOMAIN/fullchain.pem"
PRECERT_CONF="/etc/nginx/conf.d/pre-cert.conf"
DEFAULT_CONF="/etc/nginx/conf.d/default.conf"
NGINX_LOG=/var/log/nginx/error.log

# Если ENV не production — сразу запускаем nginx
if [ "$ENV" != "production" ]; then
  echo ">>> Не production-окружение. Запускаем nginx..."
  nginx -g "daemon off;"
  exit 0
fi

# Если сертификат уже есть — запускаем nginx с боевым конфигом
if [ -f "$CERT_PATH" ]; then
  echo ">>> Сертификат найден. Запускаем nginx с default.conf..."
  cp "$DEFAULT_CONF" /etc/nginx/conf.d/default.conf
  nginx -g "daemon off;"
  exit 0
fi

# Если сертификата нет — сначала запускаем nginx с pre-cert конфигом
echo ">>> Сертификат не найден. Запускаем nginx с pre-cert конфигом..."
cp "$PRECERT_CONF" /etc/nginx/conf.d/default.conf
nginx

echo ">>> Пытаемся получить сертификат..."
if certbot certonly --webroot -w /var/www/certbot \
  --email "$EMAIL" --agree-tos --no-eff-email -d "$DOMAIN"; then

  echo ">>> Сертификат успешно получен. Перезапускаем nginx с default.conf..."
  cp "$DEFAULT_CONF" /etc/nginx/conf.d/default.conf
  nginx -s reload
else
  echo ">>> Не удалось получить сертификат. Остаёмся на pre-cert.conf"
  tail -f "$NGINX_LOG"
fi

# Оставляем nginx работать в foreground
nginx -g "daemon off;"