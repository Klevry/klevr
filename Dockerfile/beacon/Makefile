docker_version="latest"

hello:
	#
	# You can use - build / run / docker / all
	#
require:
#	go get -u
	curl -sL bit.ly/go_req | bash
build: require
	go build -o ./beacon ./main.go
run:
	go run ./main.go
login:
	docker login

docker: login
	docker build -t klevry/beacon:${docker_version} ./ --no-cache
	docker push klevry/beacon:${docker_version}

all: build docker

