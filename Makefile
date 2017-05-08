ARTIFACT=len.to.tgz
HOST=alaska
DEPLOY_SCRIPT=deploy_lento
DEV_URL=http://localhost:1313

.PHONY: build dev

all: dev

deploy:
	hugo
	tar -czf $(ARTIFACT) public
	scp $(ARTIFACT) $(HOST):~/
	ssh $(HOST) ./$(DEPLOY_SCRIPT)
	rm $(ARTIFACT)

dev:
	open $(DEV_URL)
	hugo server --buildDrafts
