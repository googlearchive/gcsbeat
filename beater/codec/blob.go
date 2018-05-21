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
	"encoding/base64"

	"io"
	"io/ioutil"

	"github.com/elastic/beats/libbeat/common"
)

func NewBlobCodec(path string, input io.Reader) Codec {

	bytes, err := ioutil.ReadAll(input)
	text := ""
	if err == nil {
		text = base64.StdEncoding.EncodeToString(bytes)
	}

	event := common.MapStr{
		"event": text,
		"file":  path,
		"line":  1,
	}

	return &BlobCodec{
		event:   event,
		err:     err,
		hasMore: err == nil,
	}
}

// BufioCodec is a basic codec that reads a file line by line and reports the contents
// line number and which file the line came from
type BlobCodec struct {
	event   common.MapStr
	err     error
	hasMore bool
}

func (codec *BlobCodec) Next() bool {
	orig := codec.hasMore
	codec.hasMore = false
	return orig
}

func (codec *BlobCodec) Value() common.MapStr {
	return codec.event
}

func (codec *BlobCodec) Err() error {
	return codec.err
}
