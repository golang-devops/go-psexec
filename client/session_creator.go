package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mozillazg/request"

	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/dtos"
)

type sessionCreator struct {
	pvtKey        *rsa.PrivateKey
	baseServerUrl string
	dto           *dtos.GenTokenResponseDto
	sessionToken  []byte
	msg           *dtos.GenTokenResponseMessage
	serverPubKey  *rsa.PublicKey
}

func (s *sessionCreator) RequestToken() error {
	pubPKIXBytes, err := x509.MarshalPKIXPublicKey(&s.pvtKey.PublicKey)
	if err != nil {
		return err
	}

	c := new(http.Client)
	req := request.NewRequest(c)
	req.Json = &dtos.GetTokenRequestDto{pubPKIXBytes}

	url := combineServerUrl(s.baseServerUrl, "/token")
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Code: %d - %s", resp.StatusCode, resp.Status)
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dto := &dtos.GenTokenResponseDto{}
	err = json.Unmarshal(responseBytes, dto)
	if err != nil {
		return err
	}

	s.dto = dto
	return nil
}

func (s *sessionCreator) DecryptSessionTokenWithPrivateKey() error {
	sessionToken, err := shared.DecryptWithPrivateKey(s.pvtKey, s.dto.EncryptedSessionToken)
	if err != nil {
		return err
	}

	s.sessionToken = sessionToken
	return nil
}

func (s *sessionCreator) DecryptMessageWithSessionToken() error {
	jsonMessage, err := shared.DecryptSymmetric(s.sessionToken, s.dto.EncryptedMessage)
	if err != nil {
		return err
	}

	msg := &dtos.GenTokenResponseMessage{}
	err = json.Unmarshal(jsonMessage, msg)
	if err != nil {
		return err
	}

	s.msg = msg
	return nil
}

func (s *sessionCreator) ParseServerPublicKeyFromMessage() error {
	pubKeyInterface, err := x509.ParsePKIXPublicKey(s.msg.ServerPubKeyBytes)
	if err != nil {
		return err
	}

	serverPubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("The server public-key received is in an incorrect format")
	}

	s.serverPubKey = serverPubKey
	return nil
}

func (s *sessionCreator) VerifyServerEncryptedSessionToken() error {
	return shared.VerifySenderMessage(s.serverPubKey, s.sessionToken, s.msg.TokenEncryptionSignature)
}

func (s *sessionCreator) Create() *Session {
	return &Session{
		s.baseServerUrl,
		s.msg.SessionId,
		s.sessionToken,
		s.serverPubKey,
	}
}
