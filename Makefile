
all: build

build:
	hugo
	tar -czf len.to.tgz public
	scp len.to.tgz alaska:~/
	ssh alaska ./deploy_lento

dev:
	hugo server --buildDrafts
