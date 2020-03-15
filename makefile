.PHONY: test line build package run save load example

lint:
	golangci-lint run

test:
	go test -v

build:
	docker build -t pokemon-service .

run:
	docker run -p 3000:3000 -it pokemon-service

save:
	docker save -o image.tar pokemon-service

load:
	docker load -i image.tar

example:
	curl -X GET http://0.0.0.0:3000/pokemon/charizard
