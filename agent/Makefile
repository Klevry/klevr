hello:
	# make run/build

init:
	go mod init
require:
	go get -u

build: requirement
	go build -o ./klevr ./main.go
run:
	go run ./main.go -id=ralf -provider=baremetal -webconsole=192.168.2.100:8000
all: build
