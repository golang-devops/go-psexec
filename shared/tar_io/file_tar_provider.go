package tar_io

import (
	"fmt"
	"os"
)

func NewFileTarProvider(fullFilePath, remoteBasePath string) TarProvider {
	return &fileTarProvider{
		fullFilePath:   fullFilePath,
		remoteBasePath: remoteBasePath,
	}
}

type fileTarProvider struct {
	fullFilePath   string
	remoteBasePath string
}

func (f *fileTarProvider) IsDir() bool            { return false }
func (f *fileTarProvider) RemoteBasePath() string { return f.remoteBasePath }

func (f *fileTarProvider) Files() <-chan *TarFile {
	filesChanRW := make(chan *TarFile)

	var goRoutineErr error
	go func() {
		defer close(filesChanRW)

		//TODO: Is this correct
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
