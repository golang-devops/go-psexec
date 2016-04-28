package tar_io

import (
	"io"
	"os"
)

func NewTarFile(fileName string, contentReadCloser io.ReadCloser, isOnlyFile bool, info os.FileInfo) *TarFile {
	return &TarFile{
		FileName:          fileName,
		ContentReadCloser: contentReadCloser,
		IsOnlyFile:        isOnlyFile,
		Info:              info,
	}
}

type TarFile struct {
	FileName          string
	ContentReadCloser io.ReadCloser
	IsOnlyFile        bool
	Info              os.FileInfo
}

func (t *TarFile) HasContent() bool {
	return t.ContentReadCloser != nil
}
