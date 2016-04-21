package client

import (
	"fmt"
	"io"
)

type tarUploadHandler struct {
	session    *Session
	remotePath string
	resp       *UploadResponse
}

func (t *tarUploadHandler) Read(reader io.Reader) error {
	resp, err := t.session.UploadTarStream(t.remotePath, reader)
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
