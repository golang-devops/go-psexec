package client

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
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

func (s *Session) NewRequest() *request.Request {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Headers["Authorization"] = "Bearer " + fmt.Sprintf("sid-%d", s.sessionId)
	return req
}

func (s *Session) StreamEncryptedJsonRequest(relUrl string, rawJsonData interface{}) (*RequestResponse, error) {
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

	return &RequestResponse{Pid: int(pid), response: resp}, nil
}

func (s *Session) RequestBuilder() SessionRequestBuilderBase {
	return NewSessionRequestBuilder(s)
}

func (s *Session) EncryptAsJson(v interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return shared.EncryptSymmetric(s.sessionToken, jsonBytes)
}
