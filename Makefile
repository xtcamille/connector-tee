.PHONY: all build sign run-server run-client clean

# Variables
SERVER_DIR = tee-server
CLIENT_DIR = tee-client
SERVER_BIN = $(SERVER_DIR)/tee-server
CLIENT_BIN = $(CLIENT_DIR)/tee-client
ENCLAVE_JSON = $(SERVER_DIR)/enclave.json

# Default target
all: build

# Build the server for Occlum using standard go build
build-server:
	go build -o $(SERVER_BIN) $(SERVER_DIR)/main.go

# Build the client
build-client:
	go build -o $(CLIENT_BIN) $(CLIENT_DIR)/main.go

build: build-server build-client

# Occlum specific targets
occlum-init:
	occlum init

occlum-build: build-server
	occlum build

run-server: occlum-build
	occlum run /bin/tee-server

# Run the client
run-client:
	cd $(CLIENT_DIR) && ./tee-client

clean:
	rm -f $(SERVER_BIN) $(CLIENT_BIN)
	rm -rf occlum_instance
