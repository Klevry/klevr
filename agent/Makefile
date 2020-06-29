hello:
	# make run/build

init:
	go mod init
require:
	go get -u

build: requirement
	go build -o ./klevr ./main.go
run:
	go run ./main.go -id=ralf -provider=aws -webconsole=localhost:8000
all: build
