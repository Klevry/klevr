hello:
	# make run/build
requirement:
	curl -sL bit.ly/go_req | bash

build: requirement
	go build -o ./klevr ./main.go
run:
	go run ./main.go -id=ralf -provider=baremetal
all: build
