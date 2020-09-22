docker_version="latest"

.PHONY: help require build run login docker all

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
require: ## Check go env requirement
	go get -u
build: require ##Build Go with require process
	go build -o ./beacon ./main.go
run: ##Run main.go
	go run ./main.go
login: ##Loin in to DockerHub.io
	docker login

docker: login ##Build Docker image and Push to docker repo
	docker build -t klevry/beacon:${docker_version} ./ --no-cache
	docker push klevry/beacon:${docker_version}

all: build docker ##Build go, Docker image and push

