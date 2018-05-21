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
	"errors"
	"fmt"
	"io"

	"github.com/elastic/beats/libbeat/common"
)

const (
	JsonArrayCodecId  = "json-array"
	JsonStreamcodecId = "json-stream"
	TextCodecId       = "text"
	ClobCodecId       = "clob"
	BlobCodecId       = "blob"
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

	case codec == ClobCodecId:
		return NewClobCodec(filename, reader), nil

	case codec == BlobCodecId:
		return NewBlobCodec(filename, reader), nil

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
	return []string{
		JsonArrayCodecId,
		JsonStreamcodecId,
		TextCodecId,
		ClobCodecId,
		BlobCodecId,
	}
}
