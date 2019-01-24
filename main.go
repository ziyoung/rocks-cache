package main

import (
	"bufio"
	"io"
	"log"
	"strings"
)

// Resp is redis protocol
type Resp struct {
	Type  byte
	Value []byte
	Array []Resp
}

var nullResp = Resp{
	Type:  '$',
	Value: []byte("-1"),
	Array: nil,
}

type Reader struct {
	rd *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{bufio.NewReader(rd)}
}

func (rd *Reader) NewRespError(err error) Resp {
	return Resp{
		Type:  '-',
		Value: []byte(err.Error()),
		Array: nil,
	}
}
func (rd *Reader) ReadValue() (Resp, error) {
	b, err := rd.rd.ReadByte()
	if err != nil {
		return rd.NewRespError(err), err
	}
	resp := Resp{}
	if b == '*' {
		resp, err := rd.
	}
}

func (rd *Reader) ReadValue() (Resp, error) {

}


func main() {
	input := "*3\r\n$3\r\nset\r\n$6\r\nleader\r\n$7\r\nCharlie\r\n"
	opMap := map[byte]string{
		'+': "SimpleStrings",
		'-': "Errors",
		':': "Integers",
		'$': "BulkStrings",
		'*': "Arrays",
	}
	r := bufio.NewReader(strings.NewReader(input))
	for {
		op, err := r.ReadByte()
		if err != nil {
			log.Println(err)
			return
		}
		typ, exists := opMap[op]
		if exists {
			log.Println("op is ", typ)
		} else {
			log.Println("Unknown op ", op)
			return
		}

	}
	//log.Println("rocks-cache server start")
}
