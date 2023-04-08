export DEV=1
include .env
export 
OS:=${shell go env GOOS}
ARCH=$(shell go env GOARCH)
OOSS="linux"
ARRCHS="arm 386"
DEBUG=1
.PHONY: all clean

image:
	docker build -t auth-service:latest --build-arg OS=$(OS) --build-arg ARCH=$(ARCH)  .
run-image: image
	@docker run -p 4000:4000 -p 8083:8083 -v ./db:/db -v ./keys:/keys auth-service:latest

image-debug:
	docker build -f ./Dockerfile-debug -t auth-service:debug --build-arg OS=$(OS) --build-arg ARCH=$(ARCH)  .
run-debug: image-debug
	@docker run -p 4000:4000 -p 8083:8083 --env-file .env -v ./db:/db -v ./keys:/keys auth-service:debug

run-local:build
	@./bin/auth-$(OS)-$(ARCH)
build: clean
	@mkdir bin 
	@mkdir db
	@go build -o bin/auth-$(OS)-$(ARCH) ./main/

test:
	@go test ./...
clean:
	@rm -rf bin db