OS=linux
ARCH=x86

SERVER_BIN=bifrost-server

default: $(SERVER_BIN)

$(SERVER_BIN): ./**/*.go
	go build -o $(SERVER_BIN) cmd/server/main.go

clean:
	if [ -e $(SERVER_BIN) ]; then rm $(SERVER_BIN); fi

.PHONY: clean
