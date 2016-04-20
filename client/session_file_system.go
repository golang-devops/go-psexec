package client

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/golang-devops/go-psexec/shared"
)

type SessionFileSystem interface {
	// Download(serverUrl, localPath, remotePath string) error
	// DownloadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern string) error
	Upload(localTarReader TarReader) error
	// UploadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern string) error
	// Delete(serverUrl, remotePath string) error
	// DeleteDirFiltered(serverUrl, remotePath, dirFileFilterPattern string) error
	// Move(serverUrl, oldRemotePath, newRemotePath string) error
	// Stats(serverUrl, remotePath string) (*Stats, error)
}

func NewSessionFileSystem(session *Session) SessionFileSystem {
	return &sessionFileSystem{session: session}
}

type sessionFileSystem struct {
	session *Session
}

func (s *sessionFileSystem) checkResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		if b, e := ioutil.ReadAll(resp.Body); e != nil {
			return fmt.Errorf("The server returned status code %d but could not read response body. Error: %s", e.Error())
		} else {
			return fmt.Errorf("Server status code %d with response %s", resp.StatusCode, string(b))
		}
	}
	return nil
}

func (s *sessionFileSystem) Upload(localTarReader TarReader) error {
	pipeReader, pipeWriter := io.Pipe()
	tarWriter := tar.NewWriter(pipeWriter)
	defer pipeReader.Close() //TODO: Is this necessary?

	wg := &sync.WaitGroup{}
	wg.Add(1)

	var tarReaderError error
	go func() {
		defer wg.Done()
		defer pipeWriter.Close()
		defer tarWriter.Close()

		filesChan, errChannel := localTarReader.Files()

	forSelectLoop:
		for {
			select {
			case file, ok := <-filesChan:
				if !ok {
					break forSelectLoop
				}

				err := WriteToTar(tarWriter, file)
				if err != nil {
					tarReaderError = err
					break forSelectLoop
				}
			case errFromChan, ok := <-errChannel:
				if !ok {
					break forSelectLoop
				}
				tarReaderError = errFromChan
			}
		}

		hdr := &tar.Header{
			Name: shared.END_OF_TAR_FILENAME,
		}
		err := tarWriter.WriteHeader(hdr)
		if err != nil {
			tarReaderError = err
			return
		}
	}()

	relUrl := "/auth/upload-tar"
	remoteBasePath := localTarReader.RemoteBasePath()
	isDir := localTarReader.IsDir()
	resp, err := s.session.UploadTarStream(relUrl, remoteBasePath, isDir, pipeReader)
	if err != nil {
		return fmt.Errorf("Unable to upload tar stream, error: %s", err.Error())
	}
	defer resp.response.Body.Close()

	err = s.checkResponse(resp.response)
	if err != nil {
		return fmt.Errorf("Response error in uploading tar stream, error: %s", err.Error())
	}

	//TODO: When an error occurs right above here we do not wait, what happens with that scenario?
	wg.Wait()

	if tarReaderError != nil {
		return fmt.Errorf("Error in reading tar: %s", tarReaderError.Error())
	}

	return nil
}
