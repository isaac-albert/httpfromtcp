package main

import (
	"fmt"
	"log"
	"net"

	"www.github.com/isaac-albert/httpfromtcp/internal/request"
)

const port = 42069

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("error listening on port %v", err)
	}
	log.Printf("listening on port %d", port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting from connection: %v", err)
		}
		log.Println("connection accepted")

		go func(c net.Conn) {
			req, err := request.RequestFromReader(conn)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Request line:")
			fmt.Printf("- Method: %s\n", req.RequestLine.Method)
			fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
			fmt.Println("Headers:")
			for key, val := range req.Headers {
				fmt.Printf("- %s: %s\n", key, val)
			}
		}(conn)
	}
}


