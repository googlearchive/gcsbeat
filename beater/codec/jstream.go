package codec

import (
	"encoding/json"

	"io"

	"github.com/elastic/beats/libbeat/common"
)

func NewJsonStreamCodec(path string, input io.Reader) Codec {

	codec := &JsonStreamCodec{
		decoder:    json.NewDecoder(input),
		lineNumber: 0,
		path:       path,
	}

	return codec
}

type JsonStreamCodec struct {
	decoder    *json.Decoder
	value      common.MapStr
	err        error
	lineNumber int
	path       string
}

func (codec *JsonStreamCodec) Next() bool {
	if codec.err != nil {
		return false
	}

	jsonData := make(map[string]interface{})
	codec.err = codec.decoder.Decode(&jsonData)

	codec.lineNumber++
	codec.value = common.MapStr{
		"json": jsonData,
		"line": codec.lineNumber,
		"path": codec.path,
	}

	return codec.err == nil
}

func (codec *JsonStreamCodec) Value() common.MapStr {
	return codec.value
}

func (codec *JsonStreamCodec) Err() error {
	// EOF is not an error in our case, it's the end
	// of our stream and expected
	if codec.err == io.EOF {
		return nil
	}

	return codec.err
}
