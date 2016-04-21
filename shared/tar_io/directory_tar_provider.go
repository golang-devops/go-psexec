package tar_io

import (
	"fmt"
	"os"
	"path/filepath"
)

func NewDirectoryTarProvider(fullDirPath, filePattern, remoteBasePath string) TarProvider {
	return &directoryTarProvider{
		fullDirPath:    fullDirPath,
		filePattern:    filePattern,
		remoteBasePath: remoteBasePath,
	}
}

type directoryTarProvider struct {
	fullDirPath    string
	filePattern    string
	remoteBasePath string
}

func (d *directoryTarProvider) IsDir() bool            { return true }
func (d *directoryTarProvider) RemoteBasePath() string { return d.remoteBasePath }

func (d *directoryTarProvider) Files() <-chan *TarFile {
	filesChanRW := make(chan *TarFile)

	var goRoutineErr error
	go func() {
		defer close(filesChanRW)

		//TODO: This filepath.Walk can later be abstracted for different filesystems with `afero.Walk` from github.com/spf13/afero
		walkErr := filepath.Walk(d.fullDirPath, func(path string, info os.FileInfo, errParam error) error {
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

			relPath := path[len(d.fullDirPath):]
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
			goRoutineErr = fmt.Errorf("Unable to walk dir '%s', error: %s", d.fullDirPath, walkErr.Error())
			return
		}
	}()

	return filesChanRW
}
