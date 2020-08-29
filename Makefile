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

clone-php-client:
	if [ ! -d "./php-api-client" ]; then git clone git@github.com:obada-protocol/php-client-library ./php-api-client; fi

generate-php-client: clone-php-client
	docker run --rm \
		-v $$(pwd)/openapi:/local -v $$(pwd)/php-api-client:/src openapitools/openapi-generator-cli generate \
		-i /local/spec.openapi.yml \
		-g php \
		--skip-validate-spec \
		-o /src

export DEFAULT_NAMESPACE='OpenAPI'
export NEW_NAMESPACE='Obada'
fix-php-client-namespace:
	docker run --rm \
	  -v $$(pwd)/php-api-client:/src \
	  alpine:3.12 sh -c "grep -rl "$$DEFAULT_NAMESPACE" /src/docs | xargs sed -i 's|$$DEFAULT_NAMESPACE|$$NEW_NAMESPACE|g'" & \
	docker run --rm \
	  -v $$(pwd)/php-api-client:/src \
	  alpine:3.12 sh -c "grep -rl "$$DEFAULT_NAMESPACE" /src/lib | xargs sed -i 's|$$DEFAULT_NAMESPACE|$$NEW_NAMESPACE|g'" & \
    docker run --rm \
	  -v $$(pwd)/php-api-client:/src \
	  alpine:3.12 sh -c "grep -rl "$$DEFAULT_NAMESPACE" /src/test | xargs sed -i 's|$$DEFAULT_NAMESPACE|$$NEW_NAMESPACE|g'" & \
	docker run --rm \
	  -v $$(pwd)/php-api-client:/src \
	  alpine:3.12 sh -c "sed -i 's|$$DEFAULT_NAMESPACE|$$NEW_NAMESPACE|g' src/README.md"


deploy-php-client: generate-php-client fix-php-client-namespace
	cd php-api-client && git add . && git commit -m 'OpenApi contract update' & git push

