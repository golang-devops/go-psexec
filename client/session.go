package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/mozillazg/request"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
)

type sessionDetails struct {
	SessionId    int
	SessionToken []byte
	ServerPubKey *rsa.PublicKey
}

func (s *sessionDetails) NewRequest() *request.Request {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Headers["Authorization"] = "Bearer " + fmt.Sprintf("sid-%d", s.SessionId)
	return req
}

func (s *sessionDetails) EncryptAsJson(v interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return shared.EncryptSymmetric(s.SessionToken, jsonBytes)
}
