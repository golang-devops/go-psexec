package client

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-devops/go-psexec/shared/dtos"
	"github.com/golang-devops/go-psexec/shared/tar_io"
)

type SessionFileSystem interface {
	DownloadTar(remotePath string, options *DownloadTarOptions, tarReceiver tar_io.TarReceiver) error
	UploadTar(tarProvider tar_io.TarProvider, remotePath string, isDir bool) error
	Delete(remotePath string) error
	Move(oldRemotePath, newRemotePath string) error
	Stats(remotePath string) (*dtos.StatsDto, error)
}

func NewSessionFileSystem(session *session) SessionFileSystem {
	return &sessionFileSystem{session: session}
}

type sessionFileSystem struct {
	session *session
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

func (s *sessionFileSystem) UploadTar(tarProvider tar_io.TarProvider, remotePath string, isDir bool) error {
	uploadHandler := &tarUploadHandler{
		session:    s.session,
		remotePath: remotePath,
		isDir:      isDir,
	}

	err := tar_io.UploadProvider(tarProvider, uploadHandler)
	if err != nil {
		return fmt.Errorf("Unable to upload tar reader, error: %s", err.Error())
	}

	return nil
}

func (s *sessionFileSystem) Delete(remotePath string) error {
	relUrl := "/auth/delete"
	dto := &dtos.FsDeleteDto{Path: remotePath}
	return s.session.DoEncryptedJsonRequest(relUrl, dto, nil)
}

func (s *sessionFileSystem) Move(oldRemotePath, newRemotePath string) error {
	relUrl := "/auth/move"
	dto := &dtos.FsMoveDto{OldRemotePath: oldRemotePath, NewRemotePath: newRemotePath}
	return s.session.DoEncryptedJsonRequest(relUrl, dto, nil)
}

func (s *sessionFileSystem) Stats(remotePath string) (*dtos.StatsDto, error) {
	relUrl := "/auth/stats"
	queryValues := make(url.Values)
	queryValues.Set("path", remotePath)
	stats := &dtos.StatsDto{}
	if err := s.session.MakeGetRequest(relUrl, queryValues, stats); err != nil {
		return nil, err
	}
	return stats, nil
}
