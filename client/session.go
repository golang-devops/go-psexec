package client

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mozillazg/request"

	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/dtos"
	"github.com/golang-devops/go-psexec/shared/tar_io"
)

type Session interface {
	SessionId() int
	SessionToken() []byte
	GetFullServerUrl(relUrl string) string
	ExecRequestBuilder() SessionExecRequestBuilderBase
	FileSystem() SessionFileSystem
	EncryptAsJson(v interface{}) ([]byte, error)
}

type session struct {
	baseServerUrl string
	sessionId     int
	sessionToken  []byte
	//TODO: Currently not encrypting anything with the server public key but only session token
	serverPubKey *rsa.PublicKey
}

func (s *session) SessionId() int {
	return s.sessionId
}

func (s *session) SessionToken() []byte {
	return s.sessionToken
}

func (s *session) GetFullServerUrl(relUrl string) string {
	return combineServerUrl(s.baseServerUrl, relUrl)
}

func (s *session) newHttpClient() *http.Client {
	return new(http.Client)
}

func (s *session) NewRequest() *request.Request {
	c := s.newHttpClient()
	req := request.NewRequest(c)
	req.Headers["Authorization"] = "Bearer " + fmt.Sprintf("sid-%d", s.sessionId)
	return req
}

func (s *session) newNativeHttpRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+fmt.Sprintf("sid-%d", s.sessionId))
	return req, nil
}

func (s *session) StreamEncryptedJsonRequest(relUrl string, rawJsonData interface{}) (*ExecResponse, error) {
	url := s.GetFullServerUrl(relUrl)

	encryptedJson, err := s.EncryptAsJson(rawJsonData)
	if err != nil {
		return nil, fmt.Errorf("Unable to encrypt DTO as JSON, error: %s", err.Error())
	}

	req := s.NewRequest()
	req.Json = dtos.EncryptedJsonContainer{encryptedJson}

	resp, err := req.Post(url)
	if err != nil {
		return nil, err
	}

	pidHeader := resp.Header.Get(shared.PROCESS_ID_HTTP_HEADER_NAME)
	pid, err := strconv.ParseInt(pidHeader, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse ProcessID header '%s', error: %s", shared.PROCESS_ID_HTTP_HEADER_NAME, err.Error())
	}

	return &ExecResponse{Pid: int(pid), response: resp}, nil
}

func (s *session) DoEncryptedJsonRequest(relUrl string, inputData, destinationObject interface{}) error {
	url := s.GetFullServerUrl(relUrl)

	encryptedJson, err := s.EncryptAsJson(inputData)
	if err != nil {
		return fmt.Errorf("Unable to encrypt DTO as JSON, error: %s", err.Error())
	}

	req := s.NewRequest()
	req.Json = dtos.EncryptedJsonContainer{encryptedJson}

	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = checkResponse(resp.Response); err != nil {
		return err
	}

	if destinationObject == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(destinationObject)
}

func (s *session) MakeGetRequest(relUrl string, queryValues url.Values, destinationObject interface{}) error {
	relUrlWithQuery := strings.TrimRight(relUrl, "?") + "?" + queryValues.Encode()
	url := s.GetFullServerUrl(relUrlWithQuery)

	req, err := s.newNativeHttpRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to create http request (%s), error: %s", url, err.Error())
	}

	resp, err := s.newHttpClient().Do(req)

	if err != nil {
		return fmt.Errorf("Unable to make GET request (%s), error: %s", url, err.Error())
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return fmt.Errorf("Response error in making GET request (%s), error: %s", url, err.Error())
	}

	if destinationObject == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(destinationObject)
}

func (s *session) UploadTarStream(remotePath string, reader io.Reader) (*UploadResponse, error) {
	relUrl := "/auth/upload-tar"
	url := s.GetFullServerUrl(relUrl)

	req, err := s.newNativeHttpRequest("POST", url, reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to create http request, error: %s", err.Error())
	}

	req.Header.Add(shared.BASE_PATH_HTTP_HEADER_NAME, remotePath)

	resp, err := s.newHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to make POST request, error: %s", err.Error())
	}

	return &UploadResponse{response: resp}, nil
}

func (s *session) DownloadTarStream(queryValues url.Values, tarReceiver tar_io.TarReceiver) error {
	relUrl := "/auth/download-tar?" + queryValues.Encode()
	url := s.GetFullServerUrl(relUrl)

	req, err := s.newNativeHttpRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to create http request, error: %s", err.Error())
	}

	resp, err := s.newHttpClient().Do(req)
	if err != nil {
		return fmt.Errorf("Unable to make GET request, error: %s", err.Error())
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return fmt.Errorf("Response error in downloading tar stream, error: %s", err.Error())
	}

	err = tar_io.SaveTarToReceiver(resp.Body, tarReceiver)
	if err != nil {
		return fmt.Errorf("Unable to save response body as TAR, error: %s", err.Error())
	}

	return nil
}

func (s *session) ExecRequestBuilder() SessionExecRequestBuilderBase {
	return NewSessionExecRequestBuilderBase(s)
}

func (s *session) FileSystem() SessionFileSystem {
	return NewSessionFileSystem(s)
}

func (s *session) EncryptAsJson(v interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return shared.EncryptSymmetric(s.sessionToken, jsonBytes)
}
