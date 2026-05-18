.PHONY: build build-client build-agent clean

build: build-client build-agent

build-client:
	go build -o bin/client.exe .

build-agent:
	go build -o bin/agent.exe ./agent

clean:
	rm -rf bin/