docker_version="latest"

hello:
	#
	# You can use - build / run / docker / all
	#

init:
	go mod init github.com/ralfyang

require:
	go get -u
#	curl -sL bit.ly/go_req | bash
build: require
	go build -o ./docker_img/klevr_webconsole ./main.go
run:
	go run ./main.go -port=8000
login:
	docker login

docker: login
	docker build -t klevry/webconsole:${docker_version} ./docker_img/ --no-cache
	docker push klevry/webconsole:${docker_version}

docker-compose:
	docker-compose down
	docker-compose -f ../docker-compose.yml up -d

all: build docker docker-compose

