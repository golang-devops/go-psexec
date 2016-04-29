package tar_io

import (
	"archive/tar"
	"fmt"
	"io"
)

type writerTarReceiver struct {
	alreadyHaveFile bool
	writer          io.Writer
}

func (w *writerTarReceiver) OnEntry(tarHeader *tar.Header, tarFileReader io.Reader) error {
	if w.alreadyHaveFile {
		return fmt.Errorf("Memory TarReceiver can only process a single file")
	}
	w.alreadyHaveFile = true

	_, err := io.Copy(w.writer, tarFileReader)
	if err != nil {
		return fmt.Errorf("Unable to copy stream to buffer, error: %s", err.Error())
	}

	return nil
}
