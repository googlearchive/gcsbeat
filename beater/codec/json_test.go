// +build !integration

package codec

import (
	"github.com/elastic/beats/libbeat/common"
	"reflect"
	"strings"
	"testing"
)

func TestNewJsonArrayCodecInitValid(t *testing.T) {

	var valids = []string{
		`[]`,
		`[{}]`,
		`[{"foo":44}, {"bar":true}]`,
	}

	for _, jsonStream := range valids {
		reader := strings.NewReader(jsonStream)
		codec := NewJsonArrayCodec(reader)

		if codec.Err() != nil {
			t.Errorf("Expected initialiation to work for %q", jsonStream)
		}
	}
}

func TestNewJsonArrayCodecInitInvalid(t *testing.T) {

	var valids = []string{
		`a`,
		`"a"`,
		`{"foo":44}`,
		`][`,
		`true`,
		`33`,
	}

	for _, jsonStream := range valids {
		reader := strings.NewReader(jsonStream)
		codec := NewJsonArrayCodec(reader)

		if codec.Err() == nil {
			t.Errorf("Expected initialiation to fail for %q", jsonStream)
		}
	}
}

func TestJsonArrayNextErr(t *testing.T) {
	cases := map[string]struct {
		Json      string
		Length    int
		ExpectErr bool
	}{
		"empty array": {
			Json:      "[]",
			Length:    0,
			ExpectErr: false,
		},
		"corrupt array": {
			Json:      `[{"foo":}]`,
			Length:    0,
			ExpectErr: true,
		},
		"single array": {
			Json:      `[{"foo":33}]`,
			Length:    1,
			ExpectErr: false,
		},
		"multi array": {
			Json:      `[{"foo":33},{"foo":33},{"foo":33},{"foo":33},{"foo":33}]`,
			Length:    5,
			ExpectErr: false,
		},
		"multi array early termination": {
			Json:      `[{"foo":33},{"foo":33},`,
			Length:    2,
			ExpectErr: true,
		},
	}

	for tn, tc := range cases {
		reader := strings.NewReader(tc.Json)
		c := NewJsonArrayCodec(reader)

		if c.Err() != nil {
			t.Errorf("%q | Not expected to start with an error, got %q", tn, c.Err())
		}

		counter := 0
		for c.Next() {
			counter++
		}

		hasErr := c.Err() != nil
		if hasErr != tc.ExpectErr {
			t.Errorf("%q | Got error %d, expected? %v", tn, c.Err(), tc.ExpectErr)
		}

		if counter != tc.Length {
			t.Errorf("%q | Expected to decode %d objects, got %d", tn, tc.Length, counter)
		}

	}
}

func TestJsonArrayValue(t *testing.T) {

	const data = `[{"num":3.3, "str":"33", "bool": true}]`

	reader := strings.NewReader(data)
	c := NewJsonArrayCodec(reader)

	c.Next()
	v := c.Value()

	expected := common.MapStr{
		"num":  3.3,
		"str":  "33",
		"bool": true,
	}

	if !reflect.DeepEqual(expected, v) {
		t.Errorf("Expected %v, got %v", expected, v)
	}
}
