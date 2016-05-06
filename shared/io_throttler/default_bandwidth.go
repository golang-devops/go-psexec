package io_throttler

import (
	"github.com/efarrer/iothrottler"
)

//DefaultIOThrottlingBandwidth is the global default IO throttling bandwidth (in bytes-per-second)
var DefaultIOThrottlingBandwidth = iothrottler.BytesPerSecond * 1024 * 1024 * 10 //TODO: What is a good copy throttle speed?
