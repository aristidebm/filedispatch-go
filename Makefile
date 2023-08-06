STORE := "/mnt/axione/data/files"

BIN = "dist/filedispatch"

help:
	@echo "print help"

build:
	go build -o $(BIN)

run: build
	$(BIN) $(STORE)