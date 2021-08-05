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

deploy: quick-deploy

quick-deploy: deploy
	read -r -p "What node are you want to deploy?: " NODE; \
	ansible-playbook deployment/playbook.yml --limit $(NODE) --tags=quick-deploy

full-deploy:
	read -r -p "What node are you want to deploy?: " NODE; \
	ansible-playbook deployment/playbook.yml --limit $(NODE) --tags=full-deploy

deploy-staging:
	ansible-playbook deployment/playbook.yml --limit dev.api.obada.io

deploy-local:
	ansible-playbook deployment/playbook.yml --limit gateway.obada.local --connection=local

deploy-api-clients: deploy-node-api-libraries
	@echo "Deployment of client libraries was done"

clone-node-api-library:
	if [ ! -d "./node-api-library" ]; then git clone git@github.com:obada-foundation/node-api-library ./node-api-library; fi

clone-node-api-library-csharp: ## Clone github.com/obada-foundation/node-api-library-csharp if it does not exists
	if [ ! -d "./node-api-library-csharp" ]; then git clone git@github.com:obada-foundation/node-api-library-csharp ./node-api-library-csharp; fi

clone-node-api-library-python: ## Clone github.com/obada-foundation/node-api-library-python if it does not exists
	if [ ! -d "./node-api-library-python" ]; then git clone git@github.com:obada-foundation/node-api-library-python ./node-api-library-python; fi

generate-node-api-library-python: clone-node-api-library-python
	rm -rf $$(pwd)/node-api-library-python/*
	docker run --rm \
		-v $$(pwd)/openapi:/local -v $$(pwd)/node-api-library-python:/src openapitools/openapi-generator-cli generate \
		-i /local/spec.openapi.yml \
		-g python \
		-o /src \
		-c /local/clients/python/config.yml

generate-node-api-library-csharp: clone-node-api-library-csharp
	rm -rf $$(pwd)/node-api-library-csharp/*
	docker run --rm \
		-v $$(pwd)/openapi:/local -v $$(pwd)/node-api-library-csharp:/src openapitools/openapi-generator-cli generate \
		-i /local/spec.openapi.yml \
		-g csharp \
		-o /src \
		-c /local/clients/csharp/config.yml

generate-node-api-library: clone-node-api-library
	docker run --rm \
		-v $$(pwd)/openapi:/local -v $$(pwd)/node-api-library:/src openapitools/openapi-generator-cli generate \
		-i /local/spec.openapi.yml \
		-g php \
		-o /src \
		-c /local/clients/php/config.yml

artifacts:
	docker build -f docker/Dockerfile.artifacts --no-cache --pull -t node.bin .

build-branch:
	docker build -t $(IMAGE) -f docker/Dockerfile .

publish-branch-image:
	docker push $(IMAGE)

build-release:
	docker build -t $(RELEASE_IMAGE) -f docker/Dockerfile .

build-tag:
	docker tag $(RELEASE_IMAGE) $(TAG_IMAGE)

export NODE_LIBRARIES="node-api-library node-api-library-python node-api-library-csharp"
deploy-node-api-libraries: generate-node-api-library generate-node-api-library-csharp generate-node-api-library-python
	for library in $${NODE_LIBRARIES}; do \
  		ls -la ; \
  		echo $$library ; \
		cd $$library ; \
		git add . ; \
		HAS_CHANGES_TO_COMMIT=(`git status -s|wc -c|tr -d ' '`) ; \
		if [ "$$HAS_CHANGES_TO_COMMIT" -gt 0 ]; then \
			git commit -m 'OpenApi contract update'; \
			BRANCH=(`git branch`) ; \
			git push origin $(BRANCH) ; \
		fi ; \
		pwd ; \
		cd - ; \
	done

bpd: build-branch publish-branch-image deploy

bpdf: build-branch publish-branch-image deploy-fast

lint-openapi-spec:
	docker run \
      -v $$(pwd)/openapi:/openapi/ \
      wework/speccy lint /openapi/spec.openapi.yml

run-node:
	cd src/app/node && go run main.go

export GOPRIVATE=github.com/obada-foundation
vendor:
	cd src && go mod tidy && go mod vendor

fmt:
	cd src && go fmt ./...

test:
	cd src && go test --tags "json1 fts5 secure_delete" -v ./...

test-race: ## Runs tests with race check
	cd src && go test --tags "json1 fts5 secure_delete" -race -timeout=60s -covermode=atomic -v ./...

coverage: ## Generates and shows code coverage in a browser
	cd src && go test --tags "json1 fts5 secure_delete" ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

help: ## Show this help.
	 @IFS=$$'\n' ; \
        help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
        for help_line in $${help_lines[@]}; do \
            IFS=$$'#' ; \
            help_split=($$help_line) ; \
            help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
            help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
            printf "%-30s %s\n" $$help_command $$help_info ; \
        done