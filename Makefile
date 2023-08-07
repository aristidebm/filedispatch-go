ROOT := "/mnt/axione/data/files"

BIN = "dist/filedispatch"

help:
	@echo "print help"

build:
	go build -o $(BIN)

format:
	go fmt -x

run: build
	$(BIN) $(ROOT)