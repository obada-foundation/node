FROM php:fpm as base

RUN apt update \
	&& apt install nginx supervisor -y \
	&& docker-php-ext-install mysqli pdo pdo_mysql \
	&& rm -rf /var/www && mkdir /var/www

WORKDIR /app

FROM php:fpm as go-build

RUN apt update && apt install git gcc g++ make wget -y && \
	wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz && \
	tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz && \
	mkdir -p $GOPATH/src/github.com/arnaud-lb && \
	git clone https://github.com/arnaud-lb/php-go $GOPATH/src/github.com/arnaud-lb/php-go && \
	cd $GOPATH/src/github.com/arnaud-lb/php-go/ext && \
	phpize && ./configure && make

COPY aws-ion /aws-ion

WORKDIR /aws-ion

RUN /usr/local/go/bin/go get github.com/amzn/ion-go/ion && \
	/usr/local/go/bin/go get github.com/arnaud-lb/php-go/php-go && \
	/usr/local/go/bin/go build -o $GOPATH/src/github.com/arnaud-lb/php-go/ext/ion.so -buildmode c-shared .

FROM composer:latest AS composer

ARG APP_ENV=prod

COPY src/ /var/www

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

COPY --from=go-build /src/github.com/arnaud-lb/php-go/ext/ion.so /ext/
COPY --from=go-build /src/github.com/arnaud-lb/php-go/ext/modules/phpgo.so /ext/
COPY docker/app/etc/php7/php.ini /usr/local/etc/php
COPY docker/app/etc/nginx/nginx.conf /etc/nginx/nginx.conf
COPY docker/app/etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf
COPY docker/app/etc/supervisord.conf /etc/supervisord.conf

COPY --from=composer /var/www/ /app

STOPSIGNAL SIGQUIT

CMD ["/usr/bin/supervisord"]