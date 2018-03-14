package codec

import (
	"encoding/json"

	"io"

	"github.com/elastic/beats/libbeat/common"
)

func NewJsonStreamCodec(input io.Reader) Codec {

	codec := &JsonStreamCodec{
		decoder: json.NewDecoder(input),
	}

	return codec
}

type JsonStreamCodec struct {
	decoder *json.Decoder
	value   common.MapStr
	err     error
}

func (codec *JsonStreamCodec) Next() bool {
	if codec.err != nil {
		return false
	}

	codec.value = make(common.MapStr)
	codec.err = codec.decoder.Decode(&codec.value)

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
