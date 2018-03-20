// +build !integration

package config

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
)

type configTestCase struct {
	Name      string
	ExpectErr bool
	Config    *common.Config
}

func configure(name string, expectedFail bool, props map[string]interface{}) configTestCase {
	// default good config
	c := common.NewConfig()
	c.SetInt("interval", -1, 1)
	c.SetString("bucket_id", -1, "foo")
	c.SetBool("delete", -1, false)
	c.SetString("matches", -1, "*.log*")
	c.SetString("exclude", -1, "bak_*")
	c.SetString("metadata_key", -1, "x-goog-meta-gcsbeat")
	c.SetString("codec", -1, "text")
	c.SetBool("unpack_gzip", -1, false)
	c.SetString("processed_db_path", -1, "")

	if props == nil {
		return configTestCase{name, expectedFail, c}
	}

	// then add in the bad
	bad, _ := common.NewConfigFrom(props)
	out, _ := common.MergeConfigs(c, bad)

	return configTestCase{name, expectedFail, out}
}

func TestGetAndValidateInvalidProps(t *testing.T) {
	tests := []configTestCase{
		configure("default", false, nil),

		// intervals
		configure("zero interval", true, map[string]interface{}{"interval": 0}),
		configure("negative interval", true, map[string]interface{}{"interval": -1}),

		// buckets
		configure("missing id", true, map[string]interface{}{"bucket_id": ""}),

		// globs
		configure("good match", false, map[string]interface{}{"file_matches": "*"}),
		configure("bad match", true, map[string]interface{}{"file_matches": "[a-z"}),

		configure("good exclude", false, map[string]interface{}{"file_exclude": "*"}),
		configure("bad exclude", true, map[string]interface{}{"file_exclude": "[a-z"}),

		// metadata keys
		configure("whitesapce metadata", true, map[string]interface{}{"metadata_key": "\r\n\t "}),
		configure("empty metadata", true, map[string]interface{}{"metadata_key": ""}),
		configure("contains whitespace", false, map[string]interface{}{"metadata_key": " foo bar"}),

		// codecs
		configure("codec empty", true, map[string]interface{}{"codec": ""}),
		configure("codec unknown", true, map[string]interface{}{"codec": "foo"}),

		configure("codec text", false, map[string]interface{}{"codec": "text"}),
		configure("codec json array", false, map[string]interface{}{"codec": "json-array"}),
		configure("codec json stream", false, map[string]interface{}{"codec": "json-stream"}),
	}

	for _, testCase := range tests {
		_, err := GetAndValidateConfig(testCase.Config)

		wasErr := err != nil
		if wasErr != testCase.ExpectErr {
			t.Errorf("%q | Got error %v, expected error? %v", testCase.Name, err, testCase.ExpectErr)
		}
	}
}
