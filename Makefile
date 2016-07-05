
all: build

build:
	hugo
	tar -czf len.to.tgz public
	scp len.to.tgz alaska:~/
	ssh alaska ./deploy_lento
	rm len.to.tgz

dev:
	hugo server --buildDrafts
