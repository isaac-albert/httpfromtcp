# HTTPfromTCP

**Building HTTP/1.1 from TCP in Go**

`HTTPfromTCP` is an educational project that explores the fundamentals of web networking by implementing the HTTP/1.1 protocol directly on top of raw TCP sockets in Go. By deliberately avoiding Go’s high-level `net/http` abstractions, this project demonstrates how HTTP is structured, parsed, and transmitted at the network level.

---

## Project Overview

Modern web servers hide most protocol details behind frameworks and libraries. This project removes those abstractions and rebuilds HTTP/1.1 from the ground up on top of TCP. The goal is to gain a precise, implementation-level understanding of how browsers and servers communicate.

Key learning areas include:
- TCP connection handling
- HTTP request/response parsing
- HTTP/1.1 features such as Keep-Alive and chunked transfer encoding
- Streaming data over long-lived connections

---

## Key Features

### TCP-Based Communication
- Implements a web server directly on top of TCP sockets.
- No use of `net/http` for request parsing or response writing.
- Demonstrates how raw byte streams are converted into structured HTTP messages.

### Concurrent Request Handling
- Each incoming TCP connection is handled independently.
- Supports multiple simultaneous clients without blocking.
- Demonstrates real-world server concurrency patterns.

### HTTP/1.1 Compliance
- Implements **persistent (Keep-Alive) connections**.
- Supports **chunked transfer encoding**, enabling streaming responses.
- Compatible with modern web browsers.

---

## Technical Implementation

### Custom Request Parsing
- Manually parses:
  - Request line (method, path, version)
  - Headers (case-insensitive handling)
  - Optional request bodies
- Operates directly on raw byte streams from the TCP connection.

### Chunked Transfer Encoding
- Writes HTTP responses using chunked encoding.
- Correctly terminates chunks and sends trailers.
- Supports custom trailer headers such as `X-Content-SHA256`.

### Proxy Capabilities
- Includes a `ProxyHandler` that:
  - Forwards incoming requests to external services (e.g., `httpbin.org`)
  - Streams the upstream response back to the client
- Demonstrates bidirectional streaming over TCP.

### Routing & Status Handling
- Custom routing logic for paths such as:
  - `/video`
  - `/yourproblem`
- Handles standard HTTP status codes:
  - `200 OK`
  - `400 Bad Request`
  - `500 Internal Server Error`

---

## Project Structure

```text
.
├── cmd/
│   ├── httpserver/
│   │   └── main.go        # Full HTTP/1.1 server implementation
│   └── tcplistener/
│       └── main.go        # Minimal TCP listener for raw request logging
├── internal/
│   ├── request/           # HTTP request parsing logic
│   ├── response/          # HTTP response construction and writing
│   └── headers/           # Case-insensitive header handling and validation
```

# Getting Started

## Prerequisites
- Go 1.20 or later (recommended)

## Run the HTTP Server
- Start the server on port ```42069```:
```
go run cmd/httpserver/main.go
```
- Then open a browser and visit:
```
http://localhost:42069
```
- You can also test using ```curl```:
```
curl -v http://localhost:42069/
```

