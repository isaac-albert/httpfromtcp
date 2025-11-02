package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"www.github.com/isaac-albert/httpfromtcp/internal/headers"
)

const crlf = "\r\n"
const bufferSize = 8

type RequstState int

const (
	StateInit RequstState = iota
	StateParsingHeaders
	StateParsingBody
	StateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	State       RequstState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func NewRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
		Body: []byte{},
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
		if index >= len(p) {
			tmpBuf := make([]byte, 2*len(p))
			copy(tmpBuf, p)
			p = tmpBuf
		}
		//log.Printf("infinite loop starts here: %v", 59)
		bytesRead, err := reader.Read(p[index:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != StateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.State, bytesRead)
				}
				break
			}
			return nil, fmt.Errorf("error reading from connection")
		}
		//log.Printf("infinite loop starts here: %v", 70)
		index += bytesRead
		bytesParsed, err := req.parse(p[:index])
		if err != nil {
			return nil, err
		}
		//log.Printf("infinite loop starts here: %v", 76)

		copy(p, p[bytesParsed:])
		index -= bytesParsed
		//log.Printf("infinite loop starts here: %v", 80)

	}
	// val, _ := req.Headers.Get("content-length")
	// valInt, _ := strconv.Atoi(val)
	// if len(req.Body) < valInt {
	// 	return nil, fmt.Errorf("body smaller than content-length")
	// }

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for !r.isDone() {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case StateInit:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.State = StateParsingHeaders
		return n, nil
	case StateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = StateParsingBody
			return n, nil
		}

		return n, nil
	case StateParsingBody:
		//log.Printf("data before going to the body: '%s'", data)
		val, exists := r.Headers.Get("content-length")
		//log.Printf("value: '%s'", val)
		if !exists {
			if len(data) != 0 {
				//
				return 0, fmt.Errorf("data to be parsed: '%s' and content-legnth header: '%v'", data, r.Headers["content-length"])
			}
			r.State = StateDone
			return 0, nil
		}
		valInInteger, err := strconv.Atoi(val)
		//log.Printf("value in integer: %v", valInInteger)
		if err != nil {
			return 0, fmt.Errorf("error getting content-length")
		}
		//log.Printf("data: '%s'", data)
		if valInInteger == 0 {
			r.State = StateDone
			return 0, nil
		}
		//log.Printf("r.Body is: '%s'", r.Body)
		if len(r.Body) > valInInteger {
			return 0, fmt.Errorf("body length greater than content-length")
		}
		if len(r.Body) == valInInteger {
			r.State = StateDone
			return len(data), nil
		}
		r.Body = append(r.Body, data...)
		return len(data), nil
	case StateDone:
		return 0, fmt.Errorf("trying to parse in done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	indx := bytes.Index(data, []byte(crlf))
	if indx == -1 {
		return nil, 0, nil
	}

	reqLine, err := requestLineParsing(data[:indx])
	if err != nil {
		return nil, 0, err
	}
	return reqLine, indx + len(crlf), nil

}

func requestLineParsing(data []byte) (*RequestLine, error) {

	httpParts := strings.Split(string(data), " ")

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
