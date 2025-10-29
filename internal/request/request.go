package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const crlf = "\r\n"

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	requestLine, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	indx := bytes.Index(data, []byte(crlf))
	if indx == -1 {
		return nil, fmt.Errorf("not enough data")
	}

	return requestLineParsing(data[:indx])

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
