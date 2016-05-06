package io_throttler

import (
	"fmt"
	"io"

	"io/ioutil"

	"github.com/efarrer/iothrottler"
)

//CopyThrottled does a normal io.Copy but with throttling
func CopyThrottled(bandwidth iothrottler.Bandwidth, dest io.Writer, src io.Reader) (written int64, returnErr error) {
	pool := iothrottler.NewIOThrottlerPool(bandwidth)
	defer pool.ReleasePool()

	var readCloser io.ReadCloser
	if rc, ok := src.(io.ReadCloser); ok {
		readCloser = rc
	} else {
		readCloser = ioutil.NopCloser(readCloser)
	}

	throttledFile, err := pool.AddReader(readCloser)
	if err != nil {
		return 0, fmt.Errorf("Cannot add reader to copy throttler, error: %s", err.Error())
	}

	return io.Copy(dest, throttledFile)
}
