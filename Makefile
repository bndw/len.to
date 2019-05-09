TAG=bndw/len.to
ARTIFACT=len.to.tgz
HOST=alaska
CI_HOST=sysadm@len.to
DEPLOY_SCRIPT=deploy_lento

all: dev

.PHONY: clean
clean:
	rm -rf .build || true
	mkdir -p .build

.PHONY: build
build: clean
	hugo
	cp -R root/* .build/
	cp -R public/* .build/var/www/len.to/

.PHONY: deploy 
deploy: build
	tar -czf $(ARTIFACT) public
	scp $(ARTIFACT) $(HOST):~/
	ssh $(HOST) ./$(DEPLOY_SCRIPT)
	rm $(ARTIFACT)

.PHONY: deploy.ci
deploy.ci:
	@hugo
	@tar -czf $(ARTIFACT) public
	@scp -o "StrictHostKeyChecking=no" -P $(PORT) $(ARTIFACT) $(CI_HOST):~/
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
