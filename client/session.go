package client

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/mozillazg/request"
	"net/http"

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

func (s *Session) StartEncryptedJsonRequest(relUrl string, rawJsonData interface{}) (*request.Response, error) {
	url := s.GetFullServerUrl(relUrl)

	encryptedJson, err := s.EncryptAsJson(rawJsonData)
	if err != nil {
		return nil, fmt.Errorf("Unable to encrypt DTO as JSON, error: %s", err.Error())
	}

	req := s.NewRequest()
	req.Json = shared.EncryptedJsonContainer{encryptedJson}

	return req.Post(url)
}

func (s *Session) StartExecRequest(execDto *shared.ExecDto) (*request.Response, error) {
	relUrl := "/auth/exec"

	resp, err := s.StartEncryptedJsonRequest(relUrl, execDto)
	if err != nil {
		return nil, fmt.Errorf("Unable make POST request to url '%s', error: %s", relUrl, err.Error())
	}

	return resp, nil
}

func (s *Session) StartExecWinshellRequest(exe string, args ...string) (*request.Response, error) {
	return s.StartExecRequest(&shared.ExecDto{Executor: "winshell", Exe: exe, Args: args})
}

func (s *Session) StartExecBashRequest(exe string, args ...string) (*request.Response, error) {
	return s.StartExecRequest(&shared.ExecDto{Executor: "bash", Exe: exe, Args: args})
}

func (s *Session) EncryptAsJson(v interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return shared.EncryptSymmetric(s.sessionToken, jsonBytes)
}
