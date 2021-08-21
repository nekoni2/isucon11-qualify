.PHONY: build run bench

build:
	go build -o go/isucondition go/*.go

compose:
	sudo docker-compose up -d

jia:
	chmod +x bin/jiaapi-mock && ./bin/jiaapi-mock

run: build compose
	cd go && ./isucondition

init:
	curl -X POST http://localhost:3000/initialize