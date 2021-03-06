package filepath_summary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-devops/go-psexec/services/encoding/checksums"
)

//Service is the filepath summary service
type Service interface {
	GetFileSummary(dirPath, relPath string) (*FileSummary, error)
	GetDirSummary(dirPath string) (*DirSummary, error)
}

//New creates a new service instance
func New(checksumsSvc checksums.Service) Service {
	return &service{
		checksumsSvc: checksumsSvc,
	}
}

type service struct {
	checksumsSvc checksums.Service
}

func (s *service) GetFileSummary(dirPath, relFilePath string) (*FileSummary, error) {
	fullPath := filepath.Join(dirPath, relFilePath)
	checksumResult, err := s.checksumsSvc.FileChecksum(fullPath)
	if err != nil {
		return nil, err
	}

	fileStats, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot Stat file path '%s', error: %s", fullPath, err.Error())
	}

	return NewFileSummary(relFilePath, fileStats.ModTime(), checksumResult), nil
}

func (s *service) GetDirSummary(dirPath string) (*DirSummary, error) {
	flattenedFileSummaries := []*FileSummary{}
	walkErr := filepath.Walk(dirPath, func(path string, info os.FileInfo, innerErr error) error {
		if innerErr != nil {
			return innerErr
		}

		trimmedSourceDir := strings.Trim(dirPath, "/\\ ")
		relPath := path[len(trimmedSourceDir):]
		if relPath == "" {
			return nil
		}
		relPath = relPath[1:]

		if info.IsDir() {
			//TODO: Dirs are skipped for now?
			return nil
		}

		fileSummary, err := s.GetFileSummary(dirPath, relPath)
		if err != nil {
			return err
		}
		flattenedFileSummaries = append(flattenedFileSummaries, fileSummary)

		return nil
	})
	if walkErr != nil {
		return nil, fmt.Errorf("Cannot walk dir '%s' to get summaries, error: %s", dirPath, walkErr.Error())
	}

	baseDir := dirPath
	return NewDirSummary(baseDir, flattenedFileSummaries), nil
}
