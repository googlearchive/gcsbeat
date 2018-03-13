package codec

import (
	"bufio"
	"io"

	"github.com/elastic/beats/libbeat/common"
)

func NewBufioCodec(path string, input io.Reader) Codec {
	return &BufioCodec{
		scanner:    bufio.NewScanner(input),
		path:       path,
		lineNumber: 0,
	}
}

// BufioCodec is a basic codec that reads a file line by line and reports the contents
// line number and which file the line came from
type BufioCodec struct {
	scanner    *bufio.Scanner
	path       string
	lineNumber int
}

func (codec *BufioCodec) Next() bool {
	codec.lineNumber++

	return codec.scanner.Scan()
}

func (codec *BufioCodec) Value() common.MapStr {
	return common.MapStr{
		"event": codec.scanner.Text(),
		"file":  codec.path,
		"line":  codec.lineNumber,
	}
}

func (codec *BufioCodec) Err() error {
	return codec.scanner.Err()
}
