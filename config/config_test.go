// +build !integration

package config

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
)

func TestGetAndValidateConfigMissingRequiredFields(t *testing.T) {
	cases := map[string]struct {
		BucketId string
	}{
		"missing bucket id": {
			BucketId: "",
		},
	}

	for tn, tc := range cases {
		c := common.NewConfig()
		c.SetString("bucket_id", -1, tc.BucketId)
		_, err := GetAndValidateConfig(c)

		if err == nil {
			t.Errorf("%s: expected to fail", tn)
		}
	}
}

func TestGetAndValidateConfigInvalidConfigurations(t *testing.T) {
	cases := map[string]struct {
		Interval        int64
		BucketId        string
		Delete          bool
		Match           string
		Exclude         string
		MetadataKey     string
	}{
		"zero interval": {
			Interval:        0,
			BucketId:        "a",
			Delete:          true,
			Match:           "*",
			Exclude:         "",
			MetadataKey:     "x-goog-meta-gcsbeat",
		},
		"negative interval": {
			Interval:        0,
			BucketId:        "a",
			Delete:          true,
			Match:           "*",
			Exclude:         "",
			MetadataKey:     "x-goog-meta-gcsbeat",
		},
		"missing bucket id": {
			Interval:        1,
			BucketId:        "",
			Delete:          true,
			Match:           "*",
			Exclude:         "",
			MetadataKey:     "x-goog-meta-gcsbeat",
		},
		"bad match glob": {
			Interval:        1,
			BucketId:        "a",
			Delete:          false,
			Match:           `[a-z`,
			Exclude:         "",
			MetadataKey:     "x-goog-meta-gcsbeat",

		},
		"bad exclude glob": {
			Interval:        1,
			BucketId:        "a",
			Delete:          false,
			Match:           "",
			Exclude:         `[a-z`,
			MetadataKey:     "x-goog-meta-gcsbeat",
		},
		"whitespace only metadata key": {
			Interval:        1,
			BucketId:        "a",
			Delete:          false,
			Match:           "",
			Exclude:         `[a-z`,
			MetadataKey:     "\r\n\t ",
		},
		"empty metadata key": {
			Interval:        1,
			BucketId:        "a",
			Delete:          false,
			Match:           "",
			Exclude:         `[a-z`,
			MetadataKey:     "",
		},
	}

	for tn, tc := range cases {
		c := common.NewConfig()
		c.SetInt("interval", -1, tc.Interval)
		c.SetString("bucket_id", -1, tc.BucketId)
		c.SetBool("delete", -1, tc.Delete)
		c.SetString("file_matches", -1, tc.Match)
		c.SetString("file_exclude", -1, tc.Exclude)
		c.SetString("metadata_key", -1, tc.MetadataKey)

		_, err := GetAndValidateConfig(c)

		if err == nil {
			t.Errorf("%s: expected to fail", tn)
		}
	}
}

func TestGetAndValidateConfigValidConfiguration(t *testing.T) {
	c := common.NewConfig()
	c.SetInt("interval", -1, 1)
	c.SetString("bucket_id", -1, "foo")
	c.SetBool("delete", -1, false)
	c.SetString("matches", -1, "*.log*")
	c.SetString("exclude", -1, "bak_*")
	c.SetString("metadata_key", -1, "x-goog-meta-gcsbeat")

	_, err := GetAndValidateConfig(c)

	if err != nil {
		t.Errorf("expected to succeed but got error %v", err)
	}
}
