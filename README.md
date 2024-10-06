# Go Chat Application

## Overview

This project is a simple chat application built in Go that utilizes RSA encryption for secure messaging. The server generates a pair of RSA keys to encrypt messages sent between clients. Each client connects to the server and joins a messaging room where they can send and receive messages securely.

## Features

- **RSA Encryption**: Each client receives a public key for secure message transmission.
- **Socket Communication**: The server listens for incoming connections and facilitates communication between clients.
- **Messaging Room**: All connected clients can send and receive messages in real-time.

## Requirements

- Go 1.17 or higher
- Basic understanding of Go and socket programming
- Make

## Installation

1. Clone the repository:

```bash
$ git clone https://github.com/Gustrb/rsa-chatapp.git
$ cd go-chat-app
```

2. Generate the keys:

```bash
$ cd src/keygen
$ make
```

3. Run the server:

```bash
$ cd src/server
$ make
```

4. Connecting a Client:

```bash
$ cd src/client
$ make
```

## How does it work?

The server generates an RSA key pair (public and private keys) upon startup.
It listens for incoming client connections on a specified port.
When a client connects:
The server sends the public key to the client.
The client joins the messaging room.
Each message sent by a client is encrypted using the public key and relayed to all other clients in the room.
Clients decrypt received messages using the private key.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgements
Go Programming Language
RSA Encryption