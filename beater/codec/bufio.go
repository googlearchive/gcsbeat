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
