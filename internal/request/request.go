package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

const crlf = "\r\n"
const bufferSize = 8

type RequstState int

const (
	StateInit RequstState = iota
	StateDone
)
type Request struct {
	RequestLine RequestLine
	State RequstState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func NewRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (r *Request) isDone() bool {
	return r.State == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	p := make([]byte, bufferSize)
	req := NewRequest()

	index := 0

	for !req.isDone() {
		if len(p) >= cap(p) {
			tmpBuf := make([]byte, 2*len(p))
			n := copy(tmpBuf, p)
			p = tmpBuf[:n]
		}
		bytesRead, err := reader.Read(p[index:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = StateDone
				break
			}
			return nil, fmt.Errorf("error reading from connection")
		}
		index += bytesRead
		bytesParsed, err := req.parse(p[:index])
		if err != nil {
			return nil, err
		}
		
		tmpBuf := make([]byte, cap(p))
		copy(tmpBuf, p[bytesParsed:])
		p = tmpBuf

		index -= bytesParsed

	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case StateInit: 
		n, err := r.parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		
		r.State = StateDone
		return n, nil
	case StateDone:
		return 0, fmt.Errorf("trying to parse in done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	indx := bytes.Index(data, []byte(crlf))
	if indx == -1 {
		return 0, nil
	}
	log.Printf("data before parsing request line: %s", data[:indx])
	reqLine, err := requestLineParsing(data[:indx])
		if err != nil {
			return 0, err
		}
		r.RequestLine = *reqLine
	return indx + len(crlf), nil

}

func requestLineParsing(data []byte) (*RequestLine, error) {
	log.Printf("data: %s", data)
	httpParts := strings.Split(string(data), " ")
	log.Printf("httpParts: %v", httpParts)
	//checking for valid no of parts
	if len(httpParts) != 3 {
		return nil, fmt.Errorf("invalid no of request line parts")
	}

	//checking if method token contains all capitals
	for _, c := range httpParts[0] {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("method token is invalid")
		}
	}

	//checking if the path starts with a */*
	if httpParts[1][0] != '/' {
		return nil, fmt.Errorf("invalid target")
	}

	//checking the version
	if httpParts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid http version")
	}

	//parsing http version
	verionParts := strings.Split(httpParts[2], "/")
	if len(verionParts) != 2 {
		return nil, fmt.Errorf("invalid version")
	}

	return &RequestLine{
		HttpVersion:   verionParts[1],
		RequestTarget: httpParts[1],
		Method:        httpParts[0],
	}, nil
}
