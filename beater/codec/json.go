// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/elastic/beats/libbeat/common"
)

type JsonObject map[string]interface{}

func NewJsonArrayCodec(path string, input io.Reader) Codec {

	codec := &JsonArrayCodec{
		decoder:    json.NewDecoder(input),
		value:      common.MapStr{},
		err:        nil,
		lineNumber: 0,
		path:       path,
	}

	codec.init()

	return codec
}

// JsonArrayCodec iterates over a serialized JSON array of objects
//
// Inspiration for this decoding technique taken from go's JSON documentation
// https://golang.org/pkg/encoding/json/#Decoder.Decode
type JsonArrayCodec struct {
	decoder    *json.Decoder
	value      common.MapStr
	err        error
	lineNumber int
	path       string
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

func (codec *JsonArrayCodec) Value() common.MapStr {
	return codec.value
}

func (codec *JsonArrayCodec) Err() error {
	return codec.err
}
