package tar_io

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/golang-devops/go-psexec/shared/io_throttler"
)

type fileTarReceiver struct {
	file string
}

func (f *fileTarReceiver) OnEntry(tarHeader *tar.Header, tarFileReader io.Reader) error {
	fullDestPath := f.file

	parentDir := filepath.Dir(fullDestPath)
	err := os.MkdirAll(parentDir, os.FileMode(tarHeader.Mode))
	if err != nil {
		return fmt.Errorf("Unable to MkDirAll '%s', error: %s", parentDir, err.Error())
	}

	file, err := os.OpenFile(fullDestPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(tarHeader.Mode))
	if err != nil {
		return fmt.Errorf("Unable to open file '%s', error: %s", fullDestPath, err.Error())
	}

	defer func() {
		file.Close()
		os.Chtimes(fullDestPath, tarHeader.AccessTime, tarHeader.ModTime)
	}()

	_, err = io_throttler.CopyThrottled(io_throttler.DefaultIOThrottlingBandwidth, file, tarFileReader)
	if err != nil {
		return fmt.Errorf("Unable to copy stream to file '%s', error: %s", fullDestPath, err.Error())
	}

	return nil
}
