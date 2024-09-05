BINARY_NAME=qpc

# Default target
all: build

# Build the application
build:
	@echo "Building the application..."
	go build -o $(BINARY_NAME) main.go
	@echo "Build complete!"

# Clean the directory
clean:
	@echo "Cleaning up..."
	if [[ -f $(BINARY_NAME) ]]; then rm $(BINARY_NAME); fi

config:
	@echo "Copying a configuration file..."
	cp querypie-client.tmpl.yaml .querypie-client.yaml

.PHONY: test
test:
	go test ./...
