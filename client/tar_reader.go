package client

import (
	"fmt"
	"os"
	"path/filepath"
)

type TarReader interface {
	IsDir() bool
	RemoteBasePath() string
	Files() (<-chan *TarFile, <-chan error)
}

func NewDirTarReader(dir, filePattern, remoteBasePath string) TarReader {
	return &dirTarReader{
		dir:            dir,
		filePattern:    filePattern,
		remoteBasePath: remoteBasePath,
	}
}

type dirTarReader struct {
	dir            string
	filePattern    string
	remoteBasePath string
	files          []string
}

func (d *dirTarReader) IsDir() bool {
	return true
}
func (d *dirTarReader) RemoteBasePath() string {
	return d.remoteBasePath
}

func (d *dirTarReader) Files() (<-chan *TarFile, <-chan error) {
	filesChanRW := make(chan *TarFile)
	errorsRW := make(chan error)

	var goRoutineErr error
	go func() {
		defer close(filesChanRW)
		defer close(errorsRW)

		walkErr := filepath.Walk(d.dir, func(path string, info os.FileInfo, errParam error) error {
			if errParam != nil {
				return errParam
			}

			if d.filePattern != "" {
				if match, matchErr := filepath.Match(d.filePattern, info.Name()); matchErr != nil {
					return fmt.Errorf("File pattern match error. Pattern was '%s'. Error: %s", d.filePattern, matchErr.Error())
				} else if !match {
					return nil
				}
			}

			relPath := path[len(d.dir):]
			if relPath == "" {
				return nil
			}

			relPath = relPath[1:]

			contentReadCloser, err := os.OpenFile(path, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("Unable to read file '%s', error: %s", path, err.Error())
			}

			isOnlyFile := false
			filesChanRW <- NewTarFile(relPath, contentReadCloser, isOnlyFile, info)

			return nil
		})

		if walkErr != nil {
			goRoutineErr = fmt.Errorf("Unable to walk dir '%s', error: %s", d.dir, walkErr.Error())
			return
		}
	}()

	return filesChanRW, errorsRW
}
