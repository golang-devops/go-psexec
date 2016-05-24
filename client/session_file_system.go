package client

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/golang-devops/go-psexec/services/encoding/checksums"
	"github.com/golang-devops/go-psexec/services/filepath_summary"
	"github.com/golang-devops/go-psexec/shared/dtos"
	"github.com/golang-devops/go-psexec/shared/tar_io"
)

type SessionFileSystem interface {
	DownloadTar(remotePath string, options *DownloadTarOptions, tarReceiver tar_io.TarReceiver) error
	UploadTar(tarProvider tar_io.TarProvider, remotePath string, isDir bool) error
	Delete(remotePath string) error
	Move(oldRemotePath, newRemotePath string) error
	Copy(srcRemotePath, destRemotePath string) error
	Symlink(srcRemotePath, destRemoteSymlinkPath string) error
	Stats(remotePath string) (*dtos.StatsDto, error)
	DirSummary(remotePath string) (*filepath_summary.DirSummary, error)
	FileSummary(remotePath string) (*filepath_summary.FileSummary, error)
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

func (s *sessionFileSystem) Copy(srcRemotePath, destRemotePath string) error {
	relUrl := "/auth/copy"
	dto := &dtos.FsCopyDto{SrcRemotePath: srcRemotePath, DestRemotePath: destRemotePath}
	return s.session.DoEncryptedJsonRequest(relUrl, dto, nil)
}

func (s *sessionFileSystem) Symlink(srcRemotePath, destRemoteSymlinkPath string) error {
	relUrl := "/auth/symlink"
	dto := &dtos.FsSymlinkDto{SrcRemotePath: srcRemotePath, DestRemoteSymlinkPath: destRemoteSymlinkPath}
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

func (s *sessionFileSystem) getPathSummary(remotePath string) (*dtos.FilesystemSummaryDto, error) {
	relUrl := "/auth/path-summary"
	queryValues := make(url.Values)
	queryValues.Set("path", remotePath)
	dto := &dtos.FilesystemSummaryDto{}
	if err := s.session.MakeGetRequest(relUrl, queryValues, dto); err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *sessionFileSystem) createFileSummaryFromDTO(dto *dtos.FileSummary) (*filepath_summary.FileSummary, error) {
	checksum, err := checksums.NewChecksumResultFromHex(dto.ChecksumHex)
	if err != nil {
		return nil, err
	}

	return filepath_summary.NewFileSummary(dto.RelativePath, dto.ModTime, checksum), nil
}

func (s *sessionFileSystem) DirSummary(remotePath string) (*filepath_summary.DirSummary, error) {
	dto, err := s.getPathSummary(remotePath)
	if err != nil {
		return nil, err
	}

	flattenedFileSummaries := []*filepath_summary.FileSummary{}
	for _, f := range dto.FlattenedFiles {
		fileSummary, err := s.createFileSummaryFromDTO(f)
		if err != nil {
			return nil, err
		}
		flattenedFileSummaries = append(flattenedFileSummaries, fileSummary)
	}

	return filepath_summary.NewDirSummary(dto.BaseDir, flattenedFileSummaries), nil
}

func (s *sessionFileSystem) FileSummary(remotePath string) (*filepath_summary.FileSummary, error) {
	dto, err := s.getPathSummary(remotePath)
	if err != nil {
		return nil, err
	}

	fileSummary := dto.FlattenedFiles[0]

	//This is a hacky solution...
	//TODO: Refer to other todo item with timestamp `2016-05-09 20:57` for details on this hack

	fileSummary.RelativePath = filepath.Join(dto.BaseDir, fileSummary.RelativePath)
	return s.createFileSummaryFromDTO(fileSummary)
}
