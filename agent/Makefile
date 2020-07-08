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
	go build -o ./klevr ./main.go
run:
	go run ./main.go -id=ralf -provider=baremetal -webconsole=${WEBCONSOLE} -zone=dev -group=nexcloud
all: build
