package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"www.github.com/isaac-albert/httpfromtcp/internal/headers"
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
	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/") {
		ProxyHandler(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/video" {
		handleVideo(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, r)
		return
	}
	handler200(w, r)
}



func handler400(w *response.Writer, _ *request.Request) {
	msg := []byte(`
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
	hdrs := response.GetDefaultHeaders(len(msg))
	hdrs.ForceSet("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(hdrs)
	w.WriteBody(msg)
}

func handler500(w *response.Writer, _ *request.Request) {
	msg := []byte(`
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
	hdrs := response.GetDefaultHeaders(len(msg))
	hdrs.ForceSet("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusInternalServerError)
	w.WriteHeaders(hdrs)
	w.WriteBody(msg)
}

func handler200(w *response.Writer, _ *request.Request) {
	msg := []byte(`
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
	hdrs := response.GetDefaultHeaders(len(msg))
	hdrs.ForceSet("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(hdrs)
	w.WriteBody(msg)
}

func ProxyHandler(w *response.Writer, req *request.Request) {

	body := make([]byte, 0)
	lenOfBody := 0

	streamURLStr := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	url := fmt.Sprintf("https://httpbin.org/%s", streamURLStr)
	//log.Printf("the url is: '%s'", url)

	w.WriterState = response.StateWritingStatusLine
	w.WriteStatusLine(response.StatusOK)
	hdrs := response.GetDefaultHeaders(0)
	hdrs.ForceRemoveHeader("Content-Length")
	hdrs.ForceSet("Transfer-Coding", "chunked")
	hdrs.ForceSet("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(hdrs)

	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	const bufferSize = 1024
	buf := make([]byte, bufferSize)

	for {
		n, err := resp.Body.Read(buf)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			n, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
			//log.Printf("herre the error is: ")
			body = append(body, buf[:n]...)
			//log.Printf("if the point is here then error is not bodyappend")
			lenOfBody += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}

	_, err = w.WriteChunkedbodyDone()
	if err != nil {
		fmt.Println("error writing chunked body done")
	}

	trailers := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(body))
	trailers.ForceSet("X-Content-SHA256", sha256)
	trailers.ForceSet("X-Content-Length", fmt.Sprintf("%d", len(body)))
	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Println("Error writing trailers:", err)
		return
	}
	fmt.Println("Wrote trailers")
}

func handleVideo(w *response.Writer, _ *request.Request) {

	w.WriteStatusLine(response.StatusOK)
	fileName := "./assets/vim.mp4"

	file, err := os.ReadFile(fileName)
	if err != nil {
		handler500(w, nil)
		return
	}

	
	hdrs := response.GetDefaultHeaders(len(file))
	hdrs.ForceSet("Content-Type", "video/mp4")
	w.WriteHeaders(hdrs)
	w.WriteBody(file)
}