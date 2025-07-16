#!/bin/bash

set -e

DOMAIN=$(grep DOMAIN .env | cut -d '=' -f2)
NGINX_CONTAINER=$(docker compose ps -q nginx)

echo "🔐 Запрашиваем сертификат для $DOMAIN ..."
docker compose run --rm certbot

echo "✅ Сертификат получен, включаем HTTPS..."

# Подменяем конфиг
cp ./nginx/conf.d/https.conf ./nginx/conf.d/default.conf

echo "🔄 Перезапуск Nginx..."
docker compose restart nginx

echo "🎉 Готово! Сайт теперь доступен по https://$DOMAIN"