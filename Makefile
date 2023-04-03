export DEV=true
OS:=${shell go env GOOS}
ARCH=$(shell go env GOARCH)
OOSS="linux"
ARRCHS="arm 386"
.PHONY: all clean


image:
	@docker build . -t auth-service
run: 
	@go run github.com/maximotejeda/auth-service
	
build: clean
	@mkdir bin 
	@go build -o bin/auth-$(OS)-$(ARCH) main/main.go

test:
	@go test ./...
clean:
	@rm -rf bin db