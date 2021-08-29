image-uuidd:
	docker build -t uuidd -f Dockerfile.uuidd .

image:
	docker build -t uuid .

bin:
	mkdir -p bin

build-linux: bin
	cd src && GOOS=linux CGO_ENABLED=0 go build -o bin/uuid .

pull:
	docker pull debian:11
	docker pull golang

up:
	mkdir -p run/uuidd
	docker-compose up -d server
	docker-compose up client
