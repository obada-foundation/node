SHELL := /bin/bash
.DEFAULT_GOAL := build

build:
	@echo "Installation is done"

deploy:
	@echo "Deployment is done"

deploy-local:
	@echo "Local deployment is done"

test:
	@echo "Test is done"

deploy-api-clients: deploy-php-client
	@echo "Deployment of client libraries was done"

deploy-php-client:
	@echo "Deployment of PHP client was done"
