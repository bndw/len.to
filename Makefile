ARTIFACT=len.to.tgz
HOST=alaska
DEPLOY_SCRIPT=deploy_lento

.PHONY: build dev

all: dev

build:
	hugo
	tar -czf $(ARTIFACT) public
	scp $(ARTIFACT) $(HOST):~/
	ssh $(HOST) ./$(DEPLOY_SCRIPT)
	rm $(ARTIFACT)

dev:
	hugo server --buildDrafts
