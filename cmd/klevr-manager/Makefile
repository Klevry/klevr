docker_version="latest"

.PHONY: help init require build run login docker docker-compose all

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

init: ## Run "go mod init"
	go mod init

require: ## Check go requirement
	go get -u
#	curl -sL bit.ly/go_req | bash
build: ##require
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o ../../Dockerfile/manager/klevr-manager ./main.go

run: ## Run main.go with klevr-manager-local.yaml
	go run ./main.go -config ../../conf/klevr-manager-local.yml

login: ## Login to DockerHub.io
	docker login

docker: login ## Build docker image and push to repo
	cp -Rfv ../../conf/* ./Docker/conf/
	docker build -t klevry/klevr-manager:${docker_version} ../Dockerfile/manager/ --no-cache
	docker push klevry/klevr-manager:${docker_version}

docker-compose: ## Up docker-compose (There is no docker-compose yet.)
	docker-compose down
	docker-compose -f ../../docker-compose-demo.yml up -d

all: build docker ## Do 'build', 'docker' proccess in make

