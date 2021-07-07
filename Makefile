PROJECT = obada/node
COMMIT_BRANCH ?= develop
IMAGE = $(PROJECT):$(COMMIT_BRANCH)
RELEASE_IMAGE = $(PROJECT):master
TAG_IMAGE = $(PROJECT):$(COMMIT_TAG)

SHELL := /bin/bash
.DEFAULT_GOAL := help

lint:
	cd src && golangci-lint --config .golangci.yml run --print-issued-lines --out-format=github-actions ./...

run-local:
	docker-compose -f docker-compose.yml up -d --force-recreate

deploy-production:
	ansible-playbook deployment/playbook.yml --limit api.obada.io

deploy-staging:
	ansible-playbook deployment/playbook.yml --limit dev.api.obada.io

deploy-local:
	ansible-playbook deployment/playbook.yml --limit gateway.obada.local --connection=local

deploy-api-clients: deploy-node-api-library
	@echo "Deployment of client libraries was done"

clone-node-api-library:
	if [ ! -d "./node-api-library" ]; then git clone git@github.com:obada-foundation/node-api-library ./node-api-library; fi

generate-node-api-library: clone-node-api-library
	docker run --rm \
		-v $$(pwd)/openapi:/local -v $$(pwd)/node-api-library:/src openapitools/openapi-generator-cli generate \
		-i /local/spec.openapi.yml \
		-g php \
		-o /src \
		-c /local/clients/php/config.yml

build-branch:
	docker build -t $(IMAGE) -f docker/Dockerfile .

publish-branch-image:
	docker push $(IMAGE)

build-release:
	docker build -t $(RELEASE_IMAGE) -f docker/Dockerfile .

build-tag:
	docker tag $(RELEASE_IMAGE) $(TAG_IMAGE)

deploy-node-api-library: generate-node-api-library
	cd node-api-library ; \
	git add . ; \
	HAS_CHANGES_TO_COMMIT=(`git status -s|wc -c|tr -d ' '`) ; \
	if [ "$$HAS_CHANGES_TO_COMMIT" -gt 0 ]; then \
	  git commit -m 'OpenApi contract update'; \
	  git push origin master ; \
	fi

bpd: build-branch publish-branch-image deploy-staging

lint-openapi-spec:
	docker run \
      -v $$(pwd)/openapi:/openapi/ \
      wework/speccy lint /openapi/spec.openapi.yml

run-node:
	cd src/app/node && go run main.go

export GOPRIVATE=github.com/obada-foundation
vendor:
	cd src && go mod tidy && go mod vendor

help:
	@echo "Help here"