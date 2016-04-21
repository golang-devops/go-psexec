package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang-devops/go-psexec/shared/tar_io"
)

type SessionFileSystem interface {
	// Download(serverUrl, localPath, remotePath string) error
	UploadTar(tarProvider tar_io.TarProvider) error
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

func checkResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		if b, e := ioutil.ReadAll(resp.Body); e != nil {
			return fmt.Errorf("The server returned status code %d but could not read response body. Error: %s", e.Error())
		} else {
			return fmt.Errorf("Server status code %d with response %s", resp.StatusCode, string(b))
		}
	}
	return nil
}

type tarUploadHandler struct {
	session        *Session
	relUrl         string
	remoteBasePath string
	isDir          bool
	resp           *UploadResponse
}

func (t *tarUploadHandler) ReadPipe(pipeReader *io.PipeReader) error {
	resp, err := t.session.UploadTarStream(t.relUrl, t.remoteBasePath, t.isDir, pipeReader)
	if err != nil {
		return fmt.Errorf("Unable to write to pipe reader, error: %s", err.Error())
	}
	err = checkResponse(resp.response)
	if err != nil {
		return fmt.Errorf("Response error in uploading tar stream, error: %s", err.Error())
	}

	t.resp = resp
	return nil
}
func (t *tarUploadHandler) Done() error {
	if t.resp != nil {
		return t.resp.response.Body.Close()
	}
	return nil
}

func (s *sessionFileSystem) UploadTar(tarProvider tar_io.TarProvider) error {
	uploadHandler := &tarUploadHandler{
		session:        s.session,
		relUrl:         "/auth/upload-tar",
		remoteBasePath: tarProvider.RemoteBasePath(),
		isDir:          tarProvider.IsDir(),
	}

	err := tar_io.UploadTar(tarProvider, uploadHandler)
	if err != nil {
		return fmt.Errorf("Unable to upload tar reader, error: %s", err.Error())
	}

	return nil
}
