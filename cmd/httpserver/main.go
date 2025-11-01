package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"www.github.com/isaac-albert/httpfromtcp/internal/request"
	"www.github.com/isaac-albert/httpfromtcp/internal/response"
	"www.github.com/isaac-albert/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handle)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handle(w *response.Writer, r *request.Request) {
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		handlerFunc(w, r, response.StatusBadRequest)
		return
	case "/myproblem":
		handlerFunc(w, r, response.StatusInternalServerError)
		return
	default:
		handlerFunc(w, r, response.StatusOK)
		return
	}
}

func handlerFunc(w *response.Writer, req *request.Request, statusCode response.StatusCode) {
		msg := getHTMLBody(req)
		hdrs := response.GetDefaultHeaders(len(msg))
		hdrs.ForceSet("Content-Type", "text/html")
		w.WriteStatusLine(statusCode)
		w.WriteHeaders(hdrs)
		w.WriteBody(msg)
}

func getHTMLBody(r *request.Request) []byte {
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		return []byte(`
		<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
		`)
	case "/myproblem":
		return []byte(`
		<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
		`)
	default:
		return []byte(`
		<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
		`)

	}
}
