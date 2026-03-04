# Connector TEE with EGo

This project implements a secure TEE (Trusted Execution Environment) service using the [EGo framework](https://ego.dev/). It allows clients to safely execute JavaScript code within an Intel SGX enclave, ensuring code and data confidentiality and integrity through hardware-level encryption and remote attestation.

## Project Structure

- `api/`: contains the service data structures (`tee.go`).
- `tee-server/`: contains the TEE server implementation that runs inside the enclave.
- `tee-client/`: contains the client implementation that interacts with the TEE server and performs remote attestation.
- `Makefile`: automates common tasks like building, signing, and running the enclave.

## Prerequisites

- [Intel SGX](https://www.intel.com/content/www/us/en/architecture-and-technology/software-guard-extensions.html) enabled hardware and driver.
- [EGo](https://ego.dev/docs/getting-started/install) installed.
- [Go](https://golang.org/doc/install) installed.

## Getting Started

### 1. Build and Sign the Enclave

```bash
make build
make sign
```

### 2. Run the TEE Server

```bash
make run-server
```

### 3. Run the Client (in another terminal)

```bash
make run-client
```

## How it Works

1.  **Enclave Security**: The server is built using `ego-go` and runs within an Intel SGX enclave.
2.  **Remote Attestation**: The client uses EGo's RA-TLS (Remote Attestation TLS) to establish a secure connection with the server. It verifies that the server is indeed running the expected code inside a genuine Intel SGX enclave before sending any sensitive data.
3.  **Secure Execution**: The server receives JavaScript code and parameters over a JSON-over-TLS protocol, executes them within the enclave using the `goja` JavaScript engine, and returns the result securely to the client.
