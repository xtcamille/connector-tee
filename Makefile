.PHONY: all build sign run-server run-client clean

# Variables
SERVER_DIR = tee-server
CLIENT_DIR = tee-client
SERVER_BIN = $(SERVER_DIR)/tee-server
CLIENT_BIN = $(CLIENT_DIR)/tee-client
ENCLAVE_JSON = $(SERVER_DIR)/enclave.json

# Default target
all: build sign

# Build the server for SGX using ego-go
build-server:
	cd $(SERVER_DIR) && ego-go build -o tee-server main.go

# Build the client
build-client:
	cd $(CLIENT_DIR) && go build -o tee-client main.go

build: build-server build-client

# Sign the enclave
sign: build-server
	cd $(SERVER_DIR) && ego sign tee-server

# Run the server in enclave
run-server: sign
	cd $(SERVER_DIR) && ego run tee-server

# Run the client
run-client:
	cd $(CLIENT_DIR) && ./tee-client

clean:
	rm -f $(SERVER_BIN) $(CLIENT_BIN) $(SERVER_DIR)/tee-server.sig
