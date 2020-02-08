# NOTE: deploy.ci depends on the following environment variables:
# - CI_HOST
# - PORT
REPO ?= bndw/len.to

# Docker tags
CIRCLE_SHA1 ?= $(shell git rev-parse --short HEAD)
TAG_COMMIT=$(REPO):$(CIRCLE_SHA1)
TAG_LATEST=$(REPO):latest

# k8s deployment
GITOPS_REPO ?= git@github.com:bndw/alaska.git
GITOPS_DEPLOYMENT ?= workloads/len.to/deployment.yaml
GITOPS_CONTAINER ?= lento

# legacy deployment
ARTIFACT=len.to.tgz
HOST=alaska
DEPLOY_SCRIPT=deploy_lento

all: dev

.PHONY: build
build:
	@docker build -t $(TAG_LATEST) .

publish:
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	@docker push $(TAG_LATEST)
	@docker tag $(TAG_LATEST) $(TAG_COMMIT)
	@docker push $(TAG_COMMIT)

.PHONY: deploy.ci
deploy.ci:
	@hugo
	@cp robots.txt public/
	@tar -czf $(ARTIFACT) public
	@scp -o "StrictHostKeyChecking=no" -P $(PORT) $(ARTIFACT) $(CI_HOST):~/ > /dev/null 2>&1
	@ssh -o "StrictHostKeyChecking=no" -p $(PORT) $(CI_HOST) ./$(DEPLOY_SCRIPT)
	@rm $(ARTIFACT)

.PHONY: deploy.k8s
deploy.k8s:
	@echo "Installing kubectl"
	@curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.17.0/bin/linux/amd64/kubectl
	@chmod +x kubectl && sudo mv kubectl /usr/local/bin/
	@echo "Cloning $(GITOPS_REPO)"
	@git clone $(GITOPS_REPO) _gitopsrepo
	kubectl patch --local -o yaml -f "_gitopsrepo/$(GITOPS_DEPLOYMENT)" -p '{"spec":{"template":{"spec":{"containers":[{"name":"'$(GITOPS_CONTAINER)'","image":"'$(TAG_COMMIT)'"}]}}}}' > tmp.yaml
	@mv tmp.yaml "_gitopsrepo/$(GITOPS_DEPLOYMENT)"
	@git config --global user.email "circleci@bdw.to"
	@git config --global user.name "CircleCI"
	@cd _gitopsrepo && git add . && git commit -m "Deploy $(TAG_COMMIT)" && git push

.PHONY: dev
dev:
	open http://localhost:1313
	hugo server

.PHONY: build-rand
build-rand:
	GOOS=linux go build -o bin/rand-linux ./cmd/rand
	go build -o bin/rand ./cmd/rand
