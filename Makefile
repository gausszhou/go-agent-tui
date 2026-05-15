.PHONY: build build-client build-agent clean fmt lint vet check

build: build-client build-agent

build-client:
	go build -o bin/client.exe .

build-agent:
	go build -o bin/agent.exe ./agent/cmd

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

check: fmt vet lint

clean:
	rm -rf bin/
