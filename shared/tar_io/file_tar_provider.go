package tar_io

import (
	"fmt"
	"os"
)

type fileTarProvider struct {
	fullFilePath string
}

func (f *fileTarProvider) Files() <-chan *TarFile {
	filesChanRW := make(chan *TarFile)

	var goRoutineErr error
	go func() {
		defer close(filesChanRW)

		tarFilePath := f.fullFilePath
		fileInfo, err := os.Stat(f.fullFilePath)
		if err != nil {
			goRoutineErr = fmt.Errorf("Unable to get file info for '%s', error: %s", f.fullFilePath, err.Error())
			return
		}

		contentReadCloser, err := os.OpenFile(f.fullFilePath, os.O_RDONLY, 0)
		if err != nil {
			goRoutineErr = fmt.Errorf("Unable to read file '%s', error: %s", f.fullFilePath, err.Error())
			return
		}

		isOnlyFile := true
		filesChanRW <- NewTarFile(tarFilePath, contentReadCloser, isOnlyFile, fileInfo)
	}()

	return filesChanRW
}
