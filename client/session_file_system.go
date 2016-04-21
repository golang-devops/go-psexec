package client

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-devops/go-psexec/shared/tar_io"
)

type SessionFileSystem interface {
	DownloadTar(remotePath string, options *DownloadTarOptions, tarReceiver tar_io.TarReceiver) error
	UploadTar(tarProvider tar_io.TarProvider, remotePath string) error
	// Delete(serverUrl, remotePath string) error
	// Move(serverUrl, oldRemotePath, newRemotePath string) error
	// Stats(serverUrl, remotePath string) (*Stats, error)
}

func NewSessionFileSystem(session *Session) SessionFileSystem {
	return &sessionFileSystem{session: session}
}

type sessionFileSystem struct {
	session *Session
}

func (s *sessionFileSystem) DownloadTar(remotePath string, options *DownloadTarOptions, tarReceiver tar_io.TarReceiver) error {
	queryValues := make(url.Values)
	queryValues.Set("path", remotePath)
	if options != nil {
		if strings.TrimSpace(options.FileFilter) != "" {
			queryValues.Set("file-filter", options.FileFilter)
		}
	}

	err := s.session.DownloadTarStream(queryValues, tarReceiver)
	if err != nil {
		return fmt.Errorf("Unable to download tar, error: %s", err.Error())
	}

	return nil
}

func (s *sessionFileSystem) UploadTar(tarProvider tar_io.TarProvider, remotePath string) error {
	uploadHandler := &tarUploadHandler{
		session:    s.session,
		remotePath: remotePath,
	}

	err := tar_io.UploadProvider(tarProvider, uploadHandler)
	if err != nil {
		return fmt.Errorf("Unable to upload tar reader, error: %s", err.Error())
	}

	return nil
}
