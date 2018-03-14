package codec

import (
	"strings"
	"testing"
)

func TestJsonStreamNextErr(t *testing.T) {
	cases := map[string]struct {
		Json      string
		Length    int
		ExpectErr bool
	}{
		"empty file": {
			Json:      "",
			Length:    0,
			ExpectErr: false,
		},
		"single element": {
			Json:      `{"foo":"bar"}`,
			Length:    1,
			ExpectErr: false,
		},
		"multi element": {
			Json:      `{"foo":"bar"}{"bar":"bazz"}`,
			Length:    2,
			ExpectErr: false,
		},
		"newline delimited": {
			Json: `{"foo":"bar"}
{"bar":"bazz"}`,
			Length:    2,
			ExpectErr: false,
		},
		"second corrupt": {
			Json:      `{"foo":"bar"}aaa`,
			Length:    1,
			ExpectErr: true,
		},
		"non object": {
			Json:      `"foo":`,
			Length:    0,
			ExpectErr: true,
		},
	}

	for tn, tc := range cases {
		reader := strings.NewReader(tc.Json)
		c := NewJsonStreamCodec(reader)

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
