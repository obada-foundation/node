FROM php:fpm as base

RUN apt update \
	&& apt install nginx supervisor -y \
	&& docker-php-ext-install mysqli pdo pdo_mysql \
	&& rm -rf /var/www && mkdir /var/www

WORKDIR /app

FROM composer:latest AS composer

ARG APP_ENV=prod

COPY services/gateway/ /var/www

WORKDIR /var/www

RUN set -e \
	&& if [ "$APP_ENV" == "prod" ]; then export ARGS="--no-dev"; fi \
	&& composer install --prefer-dist --no-scripts --no-progress --no-suggest --no-interaction $ARGS \
	&& composer dump-autoload --classmap-authoritative

FROM base

RUN mkdir -p /run/php \
    && useradd --no-create-home nginx \
    && ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log

COPY docker/gateway/etc/php7/php.ini /usr/local/etc/php
COPY docker/gateway/etc/nginx/nginx.conf /etc/nginx/nginx.conf
COPY docker/gateway/etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf
COPY docker/gateway/etc/supervisord.conf /etc/supervisord.conf

COPY --from=composer:1.10.20 /var/www/ /app

RUN chown -R www-data:www-data /app

STOPSIGNAL SIGQUIT

CMD ["/usr/bin/supervisord"]