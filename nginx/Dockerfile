FROM nginx:1.27-alpine

RUN apk add --no-cache certbot bash curl

ARG ENV




COPY ./development.conf /etc/nginx/development.conf
COPY ./pre-cert.conf /etc/nginx/pre-cert.conf
COPY ./production.conf /etc/nginx/production.conf
COPY ./init-certbot.sh /init-certbot.sh

RUN chmod +x /init-certbot.sh
RUN mkdir -p /var/www/certbot

CMD ["/init-certbot.sh"]