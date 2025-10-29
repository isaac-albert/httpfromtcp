package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const addr = "localhost:42069"

func main() {
	address, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buff := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		str, err := buff.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatal(err)
		}
	}
}
