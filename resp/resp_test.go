package resp

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadResp(t *testing.T) {
	input := "*3\r\n$3\r\nset\r\n$6\r\nleader\r\n$7\r\nCharlie\r\n"
	rd := newRespReader(strings.NewReader(input))
	resp, err := rd.read()
	if err != nil {
		t.Error(err)
	}
	if resp.Type != '*' {
		t.Errorf("expected type is * but not %c", resp.Type)
	}
	if len(resp.Array) != 3 {
		t.Error("array length should be 3")
	}
	arrays := [][]byte{[]byte("set"), []byte("leader"), []byte("Charlie")}
	for i, r := range resp.Array {
		if !bytes.Equal(arrays[i], r.Value) {
			t.Errorf("expected value is %s rather than %s", arrays[i], r.Value)
		}
	}
}
