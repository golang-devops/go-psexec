package main

import "io"

type sendTarHandler struct {
	writer io.Writer
}

func (s *sendTarHandler) Read(reader io.Reader) error {
	_, err := io.Copy(s.writer, reader)
	return err
}

func (s *sendTarHandler) Done() error {
	return nil
}
