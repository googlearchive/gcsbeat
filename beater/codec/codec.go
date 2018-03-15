package codec

import (
	"errors"
	"fmt"
	"io"

	"github.com/elastic/beats/libbeat/common"
)

const (
	JsonArrayCodecId  = "json-array"
	JsonStreamcodecId = "json-stream"
	TextCodecId       = "text"
)

type Codec interface {
	// Stateful iterator inspired by bufio.Scanner

	// Next moves the cursor to the next line of input
	Next() bool

	// Gets the properties for the current line of input
	Value() common.MapStr

	// Err returns the error that caused the Codec to stop if it terminated
	// before the stream was completed.
	Err() error
}

func NewCodec(codec, filename string, reader io.Reader) (Codec, error) {
	switch {
	case codec == JsonArrayCodecId:
		return NewJsonArrayCodec(filename, reader), nil

	case codec == JsonStreamcodecId:
		return NewJsonStreamCodec(filename, reader), nil

	case codec == TextCodecId:
		return NewBufioCodec(filename, reader), nil

	// TODO could we benefit from a CSV codec?

	default:
		msg := fmt.Sprintf("No such codec: %q", codec)
		return nil, errors.New(msg)
	}
}

func IsValidCodec(codec string) bool {
	for _, k := range ValidCodecs() {
		if k == codec {
			return true
		}
	}

	return false
}

func ValidCodecs() []string {
	// generate on the fly so caller can't destructively mutate
	return []string{JsonArrayCodecId, JsonStreamcodecId, TextCodecId}
}
