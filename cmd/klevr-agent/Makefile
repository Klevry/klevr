# Make declare first
# declare -x WEBCONSOLE="192.168.2.100:8000"
# or use under variable
# WEBCONSOLE=localhost:8080
hello:
	# make run/build

init:
	go mod init
require:
	go get -u

build: require
	env GOOS=linux GOARCH=amd64 go build -o ./Dockerfile/klevr ./main.go
run:
##	go run ./main.go -id=ralf -platform=baremetal -webconsole=${WEBCONSOLE} -zone=dev -group="klevr-a-team"
	go run ./main.go -apiKey="wnskslxptmxmrPwjd123" -platform="linux_laptop" -manager="192.168.8.51:8090" -zoneId="8282"

all: build
