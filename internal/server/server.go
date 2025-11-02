package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"www.github.com/isaac-albert/httpfromtcp/internal/request"
	"www.github.com/isaac-albert/httpfromtcp/internal/response"
)



type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

// func (h *HandlerError) Write(w *response.Writer) {

// 	msg := []byte(h.Message)

// 	w.WriteStatusLine(h.StatusCode)
// 	hdrs := response.GetDefaultHeaders(len(msg))
// 	w.WriteHeaders(hdrs)
// 	w.WriteBody(msg)
// }

type Server struct {
	Listener net.Listener
	handler  Handler
	isClosed atomic.Bool
}

func NewServer() *Server {
	return &Server{}
}

func Serve(port int, handler Handler) (*Server, error) {

	server := NewServer()
	server.isClosed.Store(false)
	server.handler = handler

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}

	server.Listener = listener

	go server.Listen()

	return server, nil

}

func (s *Server) Listen() {

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				break
			}
			if errors.Is(err, io.EOF) {
				s.isClosed.Store(true)
				return
			}
			log.Fatal(err)
		}

		go func(c net.Conn) {
			s.handle(c)
		}(conn)
	}
}



func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("code is not able to parse from request")
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}
