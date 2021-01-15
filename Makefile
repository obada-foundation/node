GATEWAY_PROJECT = obada/server-gateway
QLDB_PROJECT = obada/qldb
COMMIT_BRANCH ?= dev
GATEWAY_IMAGE = $(GATEWAY_PROJECT):$(COMMIT_BRANCH)
GATEWAY_RELEASE_IMAGE = $(GATEWAY_PROJECT):master
GATEWAY_TAG_IMAGE = $(GATEWAY_PROJECT):$(COMMIT_TAG)
QLDB_IMAGE = $(QLDB_PROJECT):$(COMMIT_BRANCH)
QLDB_RELEASE_IMAGE = $(QLDB_PROJECT):master
QLDB_TAG_IMAGE = $(QLDB_PROJECT):$(COMMIT_TAG)

SHELL := /bin/bash
.DEFAULT_GOAL := help

run-local:
	docker-compose -f docker-compose.yml up -d --force-recreate

deploy-production:
	ansible-playbook deployment/playbook.yml --limit api.obada.io

deploy-staging:
	ansible-playbook deployment/playbook.yml --limit dev.api.obada.io

deploy-local:
	ansible-playbook deployment/playbook.yml --limit gateway.obada.local --connection=local

test:
	docker exec -t

deploy-api-clients: deploy-node-api-library
	@echo "Deployment of client libraries was done"

clone-node-api-library:
	if [ ! -d "./node-api-library" ]; then git clone git@github.com:obada-protocol/node-api-library ./node-api-library; fi

generate-node-api-library: clone-node-api-library
	docker run --rm \
		-v $$(pwd)/openapi:/local -v $$(pwd)/node-api-library:/src openapitools/openapi-generator-cli generate \
		-i /local/spec.openapi.yml \
		-g php \
		--skip-validate-spec \
		-o /src \
		-c /local/clients/php/config.yml

generate-javascript-client:
	docker run --rm \
		-v $$(pwd)/openapi:/local -v $$(pwd)/javascript-api-client:/src openapitools/openapi-generator-cli generate \
		-i /local/spec.openapi.yml \
		-g javascript \
		--skip-validate-spec \
		-o /src \
		-c /local/clients/javascript/config.yml

build-gateway-branch:
	docker build -t $(GATEWAY_IMAGE) -f docker/gateway/Dockerfile . --build-arg APP_ENV=dev

build-qldb-branch:
	docker build -t $(QLDB_IMAGE) -f docker/qldb/Dockerfile .

publish-branch-image-gateway:
	docker push $(GATEWAY_IMAGE)

publish-branch-image-qldb:
	docker push $(QLDB_IMAGE)

build-gateway-release:
	docker build -t $(GATEWAY_RELEASE_IMAGE) -f docker/app/Dockerfile . --build-arg APP_ENV=prod

build-qldb-release:
	docker build -t $(QLDB_RELEASE_IMAGE) -f docker/qldb/Dockerfile .

build-gateway-tag:
	docker tag $(GATEWAY_RELEASE_IMAGE) $(GATEWAY_TAG_IMAGE)

build-qldb-tag:
	docker tag $(QLDB_RELEASE_IMAGE) $(QLDB_TAG_IMAGE)

deploy-node-api-library: generate-node-api-library
	cd node-api-library && git add . && git commit -m 'OpenApi contract update' && git push origin master

bpd: build-gateway-branch publish-branch-image-gateway deploy-staging

bpdg: build-qldb-branch publish-branch-image-qldb deploy-staging

lint-openapi-spec:
	docker run \
      -v $$(pwd)/openapi:/openapi/ \
      wework/speccy lint /openapi/spec.openapi.yml

help:
	@echo "Help here"
