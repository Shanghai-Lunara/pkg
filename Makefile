.PHONY:mod

HARBOR_DOMAIN := $(shell echo ${HARBOR})
PROJECT := lunara-common

mod:
	go mod download
	go mod tidy

# echoServer
ECHO_SERVER := "$(HARBOR_DOMAIN)/$(PROJECT)/echoserver:latest"

run-echoServer:
	go run ./cmd/echoserver/main.go

build-echoServer:
	-i docker image rm $(ECHO_SERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o echo-server ./cmd/echoserver/main.go
	cp cmd/echoserver/Dockerfile . && docker build -t $(ECHO_SERVER) .
	rm -f Dockerfile && rm -f echo-server
	docker push $(ECHO_SERVER)