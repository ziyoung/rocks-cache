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

func (r Resp) String() string {
	switch r.Type {
	case '+', '-', ':', '$':
		return string(r.Value)
	case '*':
		return fmt.Sprintf("%v", r.Array)
	}
	return ""
}

type respReader struct {
	rd *bufio.Reader
}

func (rd *respReader) read() (Resp, error) {
	typ, err := rd.rd.ReadByte()
	if err != nil {
		return nullResp, err
	}
	switch typ {
	case '$':
		return rd.readString()
	case '*':
		array, err := rd.readArray()
		return Resp{Type: '*', Array: array}, err
	default:
		log.Printf("unsupported typ %c", typ)
		return nullResp, fmt.Errorf("unsupported typ %c", typ)
	}
}

func (rd *respReader) readString() (Resp, error) {
	l, err := rd.readPrefixLen()
	if err != nil {
		return nullResp, err
	}
	buf := make([]byte, l)
	buf, err = rd.readLine(l)
	if err != nil {
		return nullResp, err
	}
	return Resp{
		Type:  '$',
		Value: buf,
	}, nil
}

func (rd *respReader) readArray() ([]Resp, error) {
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

func (rd *respReader) readPrefixLen() (int, error) {
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

func (rd *respReader) readLine(l int) ([]byte, error) {
	b, err := rd.rd.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// log.Printf("len(b) = %d; l = %d", len(b), l)
	// log.Printf("b %v", b)
	if len(b) != l+2 || !bytes.HasSuffix(b, []byte("\r\n")) {
		return nil, errors.New("invalid delim or content length")
	}
	return b[:l], nil
}

// todo
// func (rd *respReader) readInt(Resp, error) {
// }

func newRespReader(rd io.Reader) *respReader {
	return &respReader{bufio.NewReader(rd)}
}
