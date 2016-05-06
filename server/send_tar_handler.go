package main

import "io"
import "github.com/golang-devops/go-psexec/shared/io_throttler"

type sendTarHandler struct {
	writer io.Writer
}

func (s *sendTarHandler) Read(reader io.Reader) error {
	_, err := io_throttler.CopyThrottled(io_throttler.DefaultIOThrottlingBandwidth, s.writer, reader)
	return err
}

func (s *sendTarHandler) Done() error {
	return nil
}
