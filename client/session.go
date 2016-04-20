package client

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mozillazg/request"

	"github.com/golang-devops/go-psexec/shared"
)

type Session struct {
	baseServerUrl string
	sessionId     int
	sessionToken  []byte
	//TODO: Currently not encrypting anything with the server public key but only session token
	serverPubKey *rsa.PublicKey
}

func (s *Session) SessionId() int {
	return s.sessionId
}

func (s *Session) SessionToken() []byte {
	return s.sessionToken
}

func (s *Session) GetFullServerUrl(relUrl string) string {
	return combineServerUrl(s.baseServerUrl, relUrl)
}

func (s *Session) newHttpClient() *http.Client {
	return new(http.Client)
}

func (s *Session) NewRequest() *request.Request {
	c := s.newHttpClient()
	req := request.NewRequest(c)
	req.Headers["Authorization"] = "Bearer " + fmt.Sprintf("sid-%d", s.sessionId)
	return req
}

func (s *Session) newNativeHttpRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+fmt.Sprintf("sid-%d", s.sessionId))
	return req, nil
}

func (s *Session) StreamEncryptedJsonRequest(relUrl string, rawJsonData interface{}) (*ExecResponse, error) {
	url := s.GetFullServerUrl(relUrl)

	encryptedJson, err := s.EncryptAsJson(rawJsonData)
	if err != nil {
		return nil, fmt.Errorf("Unable to encrypt DTO as JSON, error: %s", err.Error())
	}

	req := s.NewRequest()
	req.Json = shared.EncryptedJsonContainer{encryptedJson}

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

func (s *Session) UploadTarStream(relUrl, remoteBasePath string, isDir bool, reader io.Reader) (*UploadResponse, error) {
	url := s.GetFullServerUrl(relUrl)

	req, err := s.newNativeHttpRequest("POST", url, reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to create http request, error: %s", err.Error())
	}

	isDirHeader := "0"
	if isDir {
		isDirHeader = "1"
	}

	req.Header.Add(shared.BASE_PATH_HTTP_HEADER_NAME, remoteBasePath)
	req.Header.Add(shared.IS_DIR_HTTP_HEADER_NAME, isDirHeader)

	resp, err := s.newHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to make request, error: %s", err.Error())
	}

	return &UploadResponse{response: resp}, nil
}

func (s *Session) ExecRequestBuilder() SessionExecRequestBuilderBase {
	return NewSessionExecRequestBuilderBase(s)
}

func (s *Session) FileSystem() SessionFileSystem {
	return NewSessionFileSystem(s)
}

func (s *Session) EncryptAsJson(v interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return shared.EncryptSymmetric(s.sessionToken, jsonBytes)
}
