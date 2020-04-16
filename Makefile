NAME=http-download-speed
BINARY=build/bin/$(NAME)

.PHONY: all binary gofmt test clean run

all: binary

binary:
	go build -ldflags \
		"-X main.version=`python3 version.py`" \
	-o $(BINARY) this_module/main

gofmt:
	go fmt ./...

test:
	go test -cover -v ./...

clean:
	rm -rf build/

run: binary
	$(BINARY) --clients-num 4 --client-bitrate 10e3 --count 4 --url http://google.com
