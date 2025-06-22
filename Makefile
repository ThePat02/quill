.PHONY: all clean build test-c shared-lib

# Default target
all: build shared-lib

# Build the binary
build:
	@echo "Building Go binary..."
	@mkdir -p bin
	go build -o bin/quill cmd/quill/main.go
	@echo "Binary built: bin/quill"

# Build the shared library
shared-lib:
	@echo "Building shared library..."
	@mkdir -p bin
	CGO_ENABLED=1 go build -buildmode=c-shared -o bin/libquill.so cmd/c/main.go
	@echo "Shared library built: bin/libquill.so"
	@echo "Header file: bin/libquill.h"

# Clean build folder
clean:
	rm -rf bin
	@echo "Cleaned build folder."

# Test the C API
test-c: shared-lib
	@echo "Building C test..."
	gcc -Wall -Wextra -std=c99 \
		-I./bin \
		-L./bin \
		-o bin/test_quill \
		examples/c/example.c \
		-lquill \
		-Wl,-rpath,$$PWD/bin
	@echo "C test built: bin/test_quill"
	@echo "Running C test..."
	LD_LIBRARY_PATH=./bin ./bin/test_quill