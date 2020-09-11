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

build: require ## Build main.go
	env GOOS=linux GOARCH=amd64 go build -o ./Dockerfile/klevr ./main.go

run: ## Run main.go with some args 
##	go run ./main.go -id=ralf -platform=baremetal -webconsole=${WEBCONSOLE} -zone=dev -group="klevr-a-team"
	go run ./main.go -apiKey="wnskslxptmxmrPwjd123" -platform="linux_laptop" -manager="192.168.0.11:8090" -zoneId="8282"

all: build ## Same with build
