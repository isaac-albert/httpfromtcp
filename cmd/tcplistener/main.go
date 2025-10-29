package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
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
			lines := getLinesChannel(conn)
			for line := range lines {
				fmt.Printf("%s\n", line)
			}
		}(conn)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	b := make([]byte, 8)
	lineChan := make(chan string)

	go func(lineChan chan string) {

		defer close(lineChan)
		defer log.Println("the channel is closed")

		currentLine := ""
		for {
			n, err := f.Read(b)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("error reading from file")
			}
			indx := bytes.Index(b[:n], []byte("\n"))
			if indx == -1 {
				currentLine += string(b[:n])
			} else {
				currentLine += string(b[:indx])
				lineChan <- currentLine
				currentLine = string(b[indx+1 : n])
			}
		}
		lineChan <- currentLine
	}(lineChan)
	return lineChan
}
