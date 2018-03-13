package codec

import (
	"github.com/elastic/beats/libbeat/common"
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
