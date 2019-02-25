TAG=bndw/len.to
ARTIFACT=len.to.tgz
HOST=alaska
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
	docker build -t $(TAG) .

.PHONY: deploy 
deploy: build
	tar -czf $(ARTIFACT) public
	scp $(ARTIFACT) $(HOST):~/
	ssh $(HOST) ./$(DEPLOY_SCRIPT)
	rm $(ARTIFACT)

.PHONY: dev
dev:
	open http://localhost:1313
	hugo server

.PHONY: build-rand
build-rand:
	GOOS=linux go build -o bin/rand-linux ./cmd/rand
	go build -o bin/rand ./cmd/rand
