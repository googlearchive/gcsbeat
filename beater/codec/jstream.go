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
