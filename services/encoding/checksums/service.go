package checksums

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

//Service is the checksums service
type Service interface {
	ReaderChecksum(reader io.Reader) (*ChecksumResult, error)
	FileChecksum(filePath string) (*ChecksumResult, error)
}

//New creates a new Service instance
func New() Service {
	return &service{}
}

type service struct{}

func (s *service) ReaderChecksum(reader io.Reader) (*ChecksumResult, error) {
	h := md5.New()
	_, err := io.Copy(h, reader)
	if err != nil {
		return nil, fmt.Errorf("Cannot copy reader into md5 hasher, error: %s", err.Error())
	}

	return &ChecksumResult{h.Sum(nil)}, nil
}

func (s *service) FileChecksum(filePath string) (*ChecksumResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Cannot open file '%s' to check for checksum, error: %s", filePath, err.Error())
	}
	defer f.Close()
	return s.ReaderChecksum(f)
}
