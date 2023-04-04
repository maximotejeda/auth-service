export DEV=true
OS:=${shell go env GOOS}
ARCH=$(shell go env GOARCH)
OOSS="linux"
ARRCHS="arm 386"
DEBUG=1
.PHONY: all clean

image:
	docker build -t auth-service:latest --build-arg OS=$(OS) --build-arg ARCH=$(ARCH)  .
run: image
	@docker run -p 4000:4000 -p 8083:8083 auth-service:latest

image-debug:
	docker build -f ./Dockerfile-debug -t auth-service:debug --build-arg OS=$(OS) --build-arg ARCH=$(ARCH)  .
run-debug: image-debug
	@docker run -p 4000:4000 -p 8083:8083 auth-service:debug
	
	
build: clean
	@mkdir bin 
	@CGO_ENABLED=1 go build -o bin/auth-$(OS)-$(ARCH) main/main.go

test:
	@go test ./...
clean:
	@rm -rf bin db