package resp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// Redis Protocol specification: https://redis.io/topics/protocol
// const (
// 	SimpleStrings = '+'
// 	Errors        = '-'
// 	Integers      = ':'
// 	BulkStrings   = '$'
// 	Arrays        = '*'
// )

// Resp is Redis protocol
type Resp struct {
	Type  byte // +, -, :, $, *
	Value []byte
	Array []Resp
}

var nullResp = Resp{}

// Reader reads value
type Reader struct {
	rd *bufio.Reader
}

func (rd *Reader) read() (Resp, error) {
	typ, err := rd.rd.ReadByte()
	if err != nil {
		return nullResp, err
	}
	switch typ {
	case '$':
		return rd.readString(false)
	case '*':
		array, err := rd.readArray()
		return Resp{Type: '*', Array: array}, err
	case '+':
		return rd.readString(true)
	default:
		log.Printf("unsupported typ %c", typ)
		return nullResp, fmt.Errorf("unsupported typ %c", typ)
	}
}

func (rd *Reader) readString(isSimpleString bool) (Resp, error) {
	l, err := rd.readPrefixLen()
	if err != nil {
		return nullResp, err
	}
	buf := make([]byte, l)
	buf, err = rd.readLine(l)
	if err != nil {
		return nullResp, err
	}
	var typ byte = '$'
	if isSimpleString {
		typ = '+'
	}
	return Resp{
		Type:  typ,
		Value: buf,
	}, nil
}

func (rd *Reader) readArray() ([]Resp, error) {
	l, err := rd.readPrefixLen()
	if err != nil {
		return nil, err
	}
	buf := make([]Resp, l)
	for i := range buf {
		resp, err := rd.read()
		if err != nil {
			return nil, err
		}
		buf[i] = resp
	}
	return buf, nil
}

func (rd *Reader) readPrefixLen() (int, error) {
	tmp, err := rd.rd.ReadString('\n')
	if err != nil {
		return 0, err
	}
	// 检查是否是合法的分隔符
	if len(tmp) <= 2 || strings.HasPrefix(tmp, "\r\n") {
		return 0, errors.New("invalid delim")
	}
	l, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	return l, nil
}

func (rd *Reader) readLine(l int) ([]byte, error) {
	b, err := rd.rd.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(b) != l+2 || !bytes.HasSuffix(b, []byte("\r\n")) {
		return nil, errors.New("invalid delim or content length")
	}
	return b[:l], nil
}

// todo
// func (rd *Reader) readInt(Resp, error) {
// }

// NewReader creates a Reader
func NewReader(rd io.Reader) *Reader {
	return &Reader{bufio.NewReader(rd)}
}

// Writer write Resp response
type Writer struct {
	wr io.Writer
}

// Error builds an error
func (wr *Writer) Error(s string) {
	fmt.Fprintf(wr.wr, "-%s\r\n", s)
}

// SimpleString builds a simple string
func (wr *Writer) SimpleString(s string) {
	fmt.Fprintf(wr.wr, "+%s\r\n", s)
}

// BulkString builds a bulk string
func (wr *Writer) BulkString(s string) {
	fmt.Fprintf(wr.wr, "$%d\r\n%s\r\n", len(s), s)
}

// Null builds null
func (wr *Writer) Null() {
	wr.wr.Write([]byte("-1\r\n"))
}

// Integer builds a integer string
func (wr *Writer) Integer(v int64) {
	fmt.Fprintf(wr.wr, ":%d\r\n", v)
}

// Array builds an array
func (wr *Writer) Array(l int) {
	fmt.Fprintf(wr.wr, "*%d\r\n", l)
}

// NewWrite creates a new Write
func NewWrite(wr io.Writer) *Writer {
	return &Writer{wr: wr}
}
