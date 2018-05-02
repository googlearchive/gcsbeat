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

package storage

import (
	"github.com/elastic/beats/libbeat/logp"
)

type Filter func(filename string) (shouldFilter bool, err error)

var (
	explainLogger = logp.NewLogger("Explain")
)

// explainLogger MUST be re-initialized some time after the beat is spun up, otherwise the
// content is not written to stderr if the user uses that flag.
func initializeExplainLogger() {
	explainLogger = logp.NewLogger("Explain")
}

func FilterAndExplain(filterName string, files []string, filter Filter) (filtered []string, err error) {
	explainLogger.Debugf("Test: %s?", filterName)
	var out []string
	for _, filename := range files {

		shouldFilter, err := filter(filename)

		if err != nil {
			return out, err
		}

		if shouldFilter {
			explainLogger.Debugf(" - %q (pass)", filename)
			out = append(out, filename)
		} else {
			explainLogger.Debugf(" - %q (fail)", filename)
		}
	}

	explainLogger.Infof("Test: %s? passed %d of %d files", filterName, len(out), len(files))
	return out, nil
}

func explainFoundFiles(source string, files []string) []string {
	explainLogger.Infof("Source %q found %d files", source, len(files))

	for _, filename := range files {
		explainLogger.Debugf(" - %q", filename)
	}

	return files
}

func InvertFilter(filter Filter) Filter {
	return func(filename string) (bool, error) {
		shouldFilter, err := filter(filename)
		return !shouldFilter, err
	}
}
