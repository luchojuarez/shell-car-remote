BINARY=executable

build:
	go build -o ${BINARY} *.go

run:
	./${BINARY}

start: build run