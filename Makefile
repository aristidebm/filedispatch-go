ROOT := "/mnt/axione/data/files"

BIN = "dist/filedispatch"

help:
	@echo "print help"

build:
	go build -o $(BIN)

format:
	go fmt -x

test:
	go test -coverprofile=coverage

run: build
	$(BIN) $(ROOT)
