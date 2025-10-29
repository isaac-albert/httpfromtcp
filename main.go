package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		log.Println("error opening file")
	}
	lineChan := getLinesChannel(file)
	for line := range lineChan {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	b := make([]byte, 8)
	lineChan := make(chan string)

	go func(lineChan chan string) {
		defer close(lineChan)
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
