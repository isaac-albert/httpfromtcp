package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"
const validKeyString = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	//printing the data to check if it's parsing correctly:
	//log.Printf("h.Parse: data: %s", data)

	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return n, done, err
	}

	if index == 0 {
		n = len(crlf)
		done = true
		return n, done, nil
	}

	bytesParsed, err := h.parseHeaderString(data[:index])
	if err != nil {
		return n, done, err
	}

	n = bytesParsed + len(crlf)
	return n, done, err
}

func (h Headers) parseHeaderString(data []byte) (n int, err error) {

	//Check for valid header format field-name:

	//log.Printf("parseHeaderString: data before trimming: '%s'", data)

	data = bytes.Trim(data, " ")

	//log.Printf("parseHeaderString: data after trimming: '%s'", data)

	key, val, Exists := bytes.Cut(data, []byte(":"))
	if !Exists {
		return 0, fmt.Errorf("invalid header format")
	}

	//log.Printf("parseHeaderString: key, value after trimming: key: '%s', val: '%s'", key, val)

	val = bytes.Trim(val, " ")

	tmpKey := bytes.TrimRight(key, " ")
	if !bytes.Equal(tmpKey, key) {
		//log.Printf("data in here should be invalid colon: data: '%s'", data)
		return 0, fmt.Errorf("invalid key format")
	}

	err = validateHeaderKey(string(key))
	if err != nil {
		return 0, err
	}

	h.Set(string(key), string(val))

	return len(data), nil
}

func (h Headers) Set(key, val string) {

	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		val = strings.Join([]string{
			v,
			val,
		}, ", ")
	}
	h[key] = val

}

func (h Headers) Get(key string) string {
	return h[key]
}

func validateHeaderKey(key string) error {
	for _, c := range key {
		if !strings.ContainsAny(string(c), validKeyString) {
			return fmt.Errorf("invalid header key")
		}
	}
	return nil
}
