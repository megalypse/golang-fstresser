MAIN_PATH=./cmd
BINARY_NAME=./bin/

.PHONY: all build clean

build:
	go build -v -o ./bin ${MAIN_PATH}/...

clean:
	rm -rf ./bin/*

compile:
	chmod +x ./build.sh
	./build.sh

run-cli:
	go run ./cmd/fstresser-cli/main.go

clear-logs:
	rm -rf ./logs/*


