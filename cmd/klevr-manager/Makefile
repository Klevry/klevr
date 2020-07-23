docker_version="latest"

hello:
	#
	# You can use - init / build / run / sudo docker / all 
	#

init:
	go mod init
#	go mod init github.com/ralfyang

require:
	go get -u
#	curl -sL bit.ly/go_req | bash
build: require
	go build -o ./Docker/klevr-manager ./main.go
run:
	go run ./main.go -config ../../conf/klevr-manager-local.yml
login:
	sudo sudo docker login

docker: login
	cp -Rfv ../../conf/* ./Docker/conf/
	sudo docker build -t klevry/webconsole:${docker_version} ./Docker/ --no-cache
	sudo docker push klevry/webconsole:${docker_version}

docker-compose:
	sudo docker-compose down
	sudo docker-compose -f ../docker-compose.yml up -d

all: build docker #docker-compose

