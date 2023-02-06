MAIN_PATH=./cmd/fstresser/main.go
BINARY_NAME=./bin/fstresser

build:
	go build -o ${BINARY_NAME} ${MAIN_PATH}

compile:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ${MAIN_PATH} && \
 	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ${MAIN_PATH} && \
 	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows ${MAIN_PATH}

run: build
	./${BINARY_NAME}

