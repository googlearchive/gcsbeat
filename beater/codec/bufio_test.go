package codec

import (
	"strings"

	"testing"
)

func TestBufioCodecNextErr(t *testing.T) {
	cases := map[string]struct {
		Data      string
		Length    int
		ExpectErr bool
	}{
		"empty file": {
			Data:      ``,
			Length:    0,
			ExpectErr: false,
		},
		"single line": {
			Data:      `foo`,
			Length:    1,
			ExpectErr: false,
		},
		"two lines": {
			Data:      "foo\nbar",
			Length:    2,
			ExpectErr: false,
		},
	}

	for tn, tc := range cases {
		reader := strings.NewReader(tc.Data)
		c := NewBufioCodec("testfile", reader)

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

func TestBufioValFunctionality(t *testing.T) {
	cases := []string{"\tfoo", "bar ", "ba zz"}
	data := "\tfoo\nbar \nba zz"

	reader := strings.NewReader(data)
	codec := NewBufioCodec("testfile", reader)

	for i, evt := range cases {
		if !codec.Next() {
			t.Error("Quit too early.")
		}

		val := codec.Value()

		// lines are indexed at 1
		if val["line"] != i+1 {
			t.Errorf("Expected line to be %d, got %d", val["line"], i+1)
		}

		if val["file"] != "testfile" {
			t.Errorf("Expected file to be 'testfile', got %q", val["file"])
		}

		if val["event"] != evt {
			t.Errorf("Expected event to be %q, got %q", evt, val["event"])
		}
	}
}
