FROM nginx:1.27-alpine

# Установим certbot
RUN apk add --no-cache certbot bash curl

COPY ./${ENV}.conf /etc/nginx/conf.d/default.conf
COPY ./init-certbot.sh /init-certbot.sh

# Убедимся что скрипт исполняемый
RUN chmod +x /init-certbot.sh

# Добавим папку под challenge-файлы
RUN mkdir -p /var/www/certbot

# Запускаем скрипт вместо обычного nginx
CMD ["/init-certbot.sh"]