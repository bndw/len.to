TAG=bndw/len.to

ARTIFACT=len.to.tgz
HOST=alaska
DEPLOY_SCRIPT=deploy_lento
DEV_URL=http://localhost:5000

.PHONY: deploy dev

all: dev

clean:
	rm -rf .build || true
	mkdir -p .build

build: clean
	hugo
	cp -R root/* .build/
	cp -R public/* .build/var/www/len.to/
	docker build -t $(TAG) .

deploy:
	hugo
	tar -czf $(ARTIFACT) public
	scp $(ARTIFACT) $(HOST):~/
	ssh $(HOST) ./$(DEPLOY_SCRIPT)
	rm $(ARTIFACT)

dev: build stop
	docker run -d --name len.to -p 5000:80 $(TAG)
	open $(DEV_URL)

stop:
	docker kill len.to || true
	docker rm len.to || true
