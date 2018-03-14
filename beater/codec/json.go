package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/elastic/beats/libbeat/common"
)

type JsonObject map[string]interface{}

func NewJsonArrayCodec(input io.Reader) Codec {

	codec := &JsonArrayCodec{
		decoder: json.NewDecoder(input),
		value:   common.MapStr{},
		err:     nil,
	}

	codec.init()

	return codec
}

// JsonArrayCodec iterates over a serialized JSON array of objects
//
// Inspiration for this decoding technique taken from go's JSON documentation
// https://golang.org/pkg/encoding/json/#Decoder.Decode
type JsonArrayCodec struct {
	decoder *json.Decoder
	value   common.MapStr
	err     error
}

func (codec *JsonArrayCodec) init() {
	// Pull off the leading bracket before decoding objects
	token, err := codec.decoder.Token()
	codec.err = err

	if fmt.Sprintf("%s", token) != "[" {
		codec.err = errors.New("Invalid start token for array parsing.")
	}
}

func (codec *JsonArrayCodec) Next() bool {
	if codec.err != nil || !codec.decoder.More() {
		return false
	}

	codec.value = make(common.MapStr)
	codec.err = codec.decoder.Decode(&codec.value)

	return codec.err == nil
}

func (codec *JsonArrayCodec) Value() common.MapStr {
	return codec.value
}

func (codec *JsonArrayCodec) Err() error {
	return codec.err
}
