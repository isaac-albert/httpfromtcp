package response

import (
	"fmt"
	"io"
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
	StateWritingTrailers
	StateDone
)

type Writer struct {
	writer      io.Writer
	WriterState WriterState
}

func NewWriter(c io.Writer) *Writer {
	return &Writer{
		writer:      c,
		WriterState: StateWritingStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.WriterState != StateWritingStatusLine {
		return fmt.Errorf("error writing status line while not in State Writing Status Line")
	}

	reasonPhrase := ReasonPhrase(statusCode)

	reasonPhraseBytes := []byte(fmt.Sprintf("HTTP/1.1 %v %s\r\n", statusCode, reasonPhrase))

	_, err := w.writer.Write(reasonPhraseBytes)
	w.WriterState = StateWritingHeaders
	return err

}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != StateWritingHeaders {
		return fmt.Errorf("error writing status line while not in State Writing Headers")
	}
defer func() { w.WriterState = StateWritingBody }()
	for key, value := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(data []byte) (int, error) {

	if w.WriterState != StateWritingBody {
		return 0, fmt.Errorf("error writing status line while not in State Writing Status Line")
	}
	defer func () {w.WriterState = StateDone} ()
	n, err := w.writer.Write(data)
	if err != nil {
		return 0, fmt.Errorf("error writing to body")
	}
	_, err = w.writer.Write([]byte("\r\n"))
	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.WriterState != StateWritingBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.WriterState)
	}
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return chunkSize, nil
}

func (w *Writer) WriteChunkedbodyDone() (int, error) {
	if w.WriterState != StateWritingBody {
		return 0, fmt.Errorf("error writing body in while state: '%v'", w.WriterState)
	}
	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	w.WriterState = StateWritingTrailers
	return n, nil
}


func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.WriterState != StateWritingTrailers {
		return fmt.Errorf("error writing trailers in non-trailer state")
	}

	defer func() {w.WriterState = StateWritingBody} ()
	
	//log.Printf("headers in trailers: %v", h)
	for key, value := range h {
		data := []byte(fmt.Sprintf("%s: %s\r\n", key, value))
		_, err := w.writer.Write(data)
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	return err
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
