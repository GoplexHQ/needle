package internal

import (
	"bytes"
	"runtime"
)

const stackBufferSize = 64

// GetGoroutineID returns the ID of the current goroutine.
func GetGoroutineID() string {
	b := make([]byte, stackBufferSize)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]

	return string(b)
}
