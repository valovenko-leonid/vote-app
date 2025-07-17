FROM nginx:1.27-alpine

RUN apk add --no-cache certbot bash curl

ARG ENV
COPY ./pre-cert.conf /etc/nginx/conf.d/default.conf
COPY ./pre-cert.conf /etc/nginx/conf.d/pre-cert.conf
COPY ./production.conf /etc/nginx/conf.d/production.conf
COPY ./init-certbot.sh /init-certbot.sh

RUN chmod +x /init-certbot.sh
RUN mkdir -p /var/www/certbot

CMD ["/init-certbot.sh"]