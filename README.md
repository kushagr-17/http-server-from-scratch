# http-server-from-scratch
Aim of this project is to create a HTTP server from scratch using Golang. Completely forbidding the usage of Go's inbuilt `net/http` package.
This focuses mainly on backend and networking

## Basic stuff
Everybody knows how a valid HTTP request looks like, the three most important pieces are:
- Method (`GET/POST/PUT` etc)
- Target Path
- Protocol Version (here, `HTTP/1.1`)
Everything after it (Host, User-Agent, Accept) are called HTTP headers which are just key-value pairs providing metadata.

A valid HTTP response has 4 essential parts:
- **Status Line**:  `HTTP/1.1 200 OK\r\n`
- **Headers**: `Content-Type: text/html\r\n`
- **Blank line**: `\r\n` telling browser that headers over, body starting
- **Body**: *self-explanatory*

## Architecture
This project implements a functional HTTP/1.1 server from scratch which interacts directly with the transport layer using raw TCP sockets via the `net` package.

### 1. TCP Listener (Transport Layer)
* The server binds to a local port using `net.Listen("tcp", ":8080")`, requesting a network socket from the operating system. 
* It enters an infinite loop, blocking execution at `listener.Accept()` until the OS passes a fully established TCP connection (the 3-way handshake is handled by the OS kernel) to the application layer.

### 2. HTTP Parser (Application Layer)
* Upon receiving a connection, the server reads the raw byte stream into a buffer.
* Web browsers communicate using the strict, text-based HTTP/1.1 protocol. The server parses this raw string by splitting it along newline characters (`\n`).
* It isolates the **Request Line** (e.g., `GET / HTTP/1.1`) and tokenizes it by spaces to extract the HTTP Method (`GET`) and the URL Path (`/`).

### 3. The Responder
* Browsers will drop connections if the response does not adhere strictly to the HTTP RFC specifications. 
* The server manually constructs the HTTP response string, ensuring the inclusion of the required Status Line (`HTTP/1.1 200 OK`), essential headers (`Content-Type`, `Content-Length`), and the mandatory Carriage Return + Line Feed (`\r\n\r\n`) that separates the headers from the HTML body.
* The formatted string is written back across the TCP socket (`conn.Write`), and the connection is closed (`conn.Close()`), simulating the stateless nature of standard HTTP.
