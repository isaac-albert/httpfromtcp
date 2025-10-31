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
	case "/":
		msg := []byte("Your request was an absolute banger.")
		hdrs := response.GetDefaultHeaders(len(msg))
		hdrs.ForceSet("Content-Type", "text/html")
		status := response.StatusOK
		// title := fmt.Sprintf("%v %s", status, response.ReasonPhrase(status))
		// heading := response.GetHeading(status)
		w.WriteStatusLine(status)
		w.WriteHeaders(hdrs)
		// 		data := []byte(fmt.Sprintf(`
		// 		<html>
		//   <head>
		//     <title>%s</title>
		//   </head>
		//   <body>
		//     <h1>%s</h1>
		//     <p>%s</p>
		//   </body>
		// </html>
		// 		`, title, heading, string(msg)))
		w.WriteBody(msg)
	case "/yourproblem":
		msg := []byte("Your request honestly kinda sucked.")
		hdrs := response.GetDefaultHeaders(len(msg))
		hdrs.ForceSet("Content-Type", "text/html")
		status := response.StatusBadRequest
		// title := fmt.Sprintf("%v %s", status, response.ReasonPhrase(status))
		// heading := response.GetHeading(status)
		w.WriteStatusLine(status)
		w.WriteHeaders(hdrs)
		// 		data := []byte(fmt.Sprintf(`
		// 		<html>
		//   <head>
		//     <title>%s</title>
		//   </head>
		//   <body>
		//     <h1>%s</h1>
		//     <p>%s</p>
		//   </body>
		// </html>
		// 		`, title, heading, string(msg)))
		w.WriteBody(msg)
	case "/myproblem":
		msg := []byte("Okay, you know what? This one is on me.")
		hdrs := response.GetDefaultHeaders(len(msg))
		hdrs.ForceSet("Content-Type", "text/html")
		status := response.StatusInternalServerError
		// title := fmt.Sprintf("%v %s", status, response.ReasonPhrase(status))
		// heading := response.GetHeading(status)
		w.WriteStatusLine(status)
		w.WriteHeaders(hdrs)
		// 		data := []byte(fmt.Sprintf(`
		// 		<html>
		//   <head>
		//     <title>%s</title>
		//   </head>
		//   <body>
		//     <h1>%s</h1>
		//     <p>%s</p>
		//   </body>
		// </html>
		// 		`, title, heading, string(msg)))
		w.WriteBody(msg)
	}
}
