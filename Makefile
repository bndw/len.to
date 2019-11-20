# NOTE: deploy.ci depends on the following environment variables:
# - CI_HOST
# - PORT
REPO=bndw/len.to
CIRCLE_SHA1 ?= $(shell git rev-parse --short HEAD)
ARTIFACT=len.to.tgz
HOST=alaska
DEPLOY_SCRIPT=deploy_lento

all: dev

.PHONY: build
build:
	@docker build -t $(REPO):latest .

publish:
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	@docker push $(REPO):latest
	@docker tag $(REPO):latest $(REPO):$(CIRCLE_SHA1)
	@docker push $(REPO):$(CIRCLE_SHA1)

.PHONY: deploy.ci
deploy.ci:
	@hugo
	@tar -czf $(ARTIFACT) public
	@scp -o "StrictHostKeyChecking=no" -P $(PORT) $(ARTIFACT) $(CI_HOST):~/ > /dev/null 2>&1
	@ssh -o "StrictHostKeyChecking=no" -p $(PORT) $(CI_HOST) ./$(DEPLOY_SCRIPT)
	@rm $(ARTIFACT)

.PHONY: dev
dev:
	open http://localhost:1313
	hugo server

.PHONY: build-rand
build-rand:
	GOOS=linux go build -o bin/rand-linux ./cmd/rand
	go build -o bin/rand ./cmd/rand
