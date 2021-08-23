.PHONY: build run bench

build:
	go build -o go/isucondition go/*.go

compose-restart:
	sudo docker-compose restart

compose:
	sudo docker-compose up -d

jia:
	chmod +x bin/jiaapi-mock && ./bin/jiaapi-mock

run: build compose
	cd go && ./isucondition

init:
	curl -X POST http://localhost:3000/initialize

# Patch for ignore self-signed TLS cert
setup-benchmark:
	rm -rf isucon-official
	git clone https://github.com/isucon/isucon11-qualify.git isucon-official --depth=1
	cp patch/bench/main.go isucon-official/bench/
	cd isucon-official/bench && go build -o bench main.go
	wget https://github.com/isucon/isucon11-qualify/releases/download/public/initialize.json -O isucon-official/bench/data/initialize.json

copy:
	scp go/isucondition isucon01:/home/isucon/webapp/go/isucondition

benchmark:
	cd isucon-official/bench && ./bench -all-addresses 127.0.0.1 -jia-service-url http://127.0.0.1:5000 -target 127.0.0.1:443 -tls true -exit-status