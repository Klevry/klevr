# Make declare first
# declare -x WEBCONSOLE="192.168.2.100:8000"
# or use under variable
# WEBCONSOLE=localhost:8080

.PHONY: help init require build run all

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

init: ## Run "go mod init"
	go mod init

require: ## Check go requirement
	go get -u

build: ##require ## Build main.go
	env GOOS=linux GOARCH=amd64 go build -o ../../Dockerfile/agent/klevr ./main.go

docker: build
	docker stop -t 0 agent || true
	docker rm agent || true
	docker rmi agent
	docker build -t agent ../../Dockerfile/agent
	docker run --name agent --add-host=host-docker-internal:172.17.0.1 -d -e K_API_KEY="751ad8ec4db84c889f4d58bd22be3953" -e K_PLATFORM="kubernetes" -e K_ZONE_ID=1 -e K_MANAGER_URL="http://192.168.15.150:8090" agent

run: ## Run main.go with some args 
	go run ./main.go -apiKey="2d732d48059f463aace751a723d7a239" -platform="kubernetes" -manager="http://localhost:8090" -zoneId="18"

all: build ## Same with build
