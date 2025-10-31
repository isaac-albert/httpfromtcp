package response

import (
	"fmt"
	"io"
	"net"
	"strconv"

	"www.github.com/isaac-albert/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
	StatusUnknown             StatusCode = 0
)

type WriterState int

const (
	StateWritingStatusLine WriterState = iota
	StateWritingHeaders
	StateWritingBody
	StateDone
)

type Writer struct {
	writer      io.Writer
	writerState WriterState
}

func NewWriter(c net.Conn) *Writer {
	return &Writer{
		writer:      c,
		writerState: StateWritingStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	reasonPhrase := ReasonPhrase(statusCode)

	reasonPhraseBytes := []byte(fmt.Sprintf("HTTP/1.1 %v %s\r\n", statusCode, reasonPhrase))

	_, err := w.writer.Write(reasonPhraseBytes)
	w.writerState = StateWritingHeaders
	return err

}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	for key, value := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	w.writerState = StateWritingBody
	return err
}

func (w *Writer) WriteBody(data []byte) (int, error) {

	n, err := w.writer.Write(data)
	return n, err
}

func GetHeading(status StatusCode) string {
	switch status {
	case StatusOK:
		return "Success!"
	case StatusBadRequest:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	contentLenInt := strconv.Itoa(contentLen)

	h.Set("Content-Length", contentLenInt)
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func ReasonPhrase(s StatusCode) string {

	switch s {
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}
