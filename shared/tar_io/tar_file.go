package tar_io

import (
	"io"
	"os"
)

type TarFile struct {
	FileName          string
	ContentReadCloser io.ReadCloser
	IsOnlyFile        bool
	Info              os.FileInfo
}

func NewTarFile(fileName string, contentReadCloser io.ReadCloser, isOnlyFile bool, info os.FileInfo) *TarFile {
	return &TarFile{
		FileName:          fileName,
		ContentReadCloser: contentReadCloser,
		IsOnlyFile:        isOnlyFile,
		Info:              info,
	}
}
