package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mozillazg/request"
	"io/ioutil"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
)

type sessionDetails struct {
	SessionId    int
	SessionToken []byte
	ServerPubKey *rsa.PublicKey
}

type sessionCreator struct {
	pvtKey       *rsa.PrivateKey
	dto          *shared.GenTokenResponseDto
	sessionToken []byte
	msg          *shared.GenTokenResponseMessage
	serverPubKey *rsa.PublicKey
}

func (s *sessionCreator) requestToken() error {
	pubPKIXBytes, err := x509.MarshalPKIXPublicKey(&s.pvtKey.PublicKey)
	if err != nil {
		return err
	}

	c := new(http.Client)
	req := request.NewRequest(c)
	req.Json = &shared.GetTokenRequestDto{pubPKIXBytes}

	url := combineServerUrl("/token")
	resp, err := req.Post(url)
	checkError(err)
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dto := &shared.GenTokenResponseDto{}
	err = json.Unmarshal(responseBytes, dto)
	if err != nil {
		return err
	}

	s.dto = dto
	return nil
}

func (s *sessionCreator) decryptSessionTokenWithPrivateKey() error {
	sessionToken, err := shared.DecryptWithPrivateKey(s.pvtKey, s.dto.EncryptedSessionToken)
	if err != nil {
		return err
	}

	s.sessionToken = sessionToken
	return nil
}

func (s *sessionCreator) decryptMessageWithSessionToken() error {
	jsonMessage, err := shared.DecryptSymmetric(s.sessionToken, s.dto.EncryptedMessage)
	if err != nil {
		return err
	}

	msg := &shared.GenTokenResponseMessage{}
	err = json.Unmarshal(jsonMessage, msg)
	if err != nil {
		return err
	}

	s.msg = msg
	return nil
}

func (s *sessionCreator) parseServerPublicKeyFromMessage() error {
	pubKeyInterface, err := x509.ParsePKIXPublicKey(s.msg.ServerPubKeyBytes)
	if err != nil {
		return err
	}

	serverPubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("The server public-key received is in an incorrect format")
	}

	s.serverPubKey = serverPubKey
	return nil
}

func (s *sessionCreator) verifyServerEncryptedSessionToken() error {
	return shared.VerifySenderMessage(s.serverPubKey, s.dto.EncryptedSessionToken, s.msg.TokenEncryptionSignature)
}

func (s *sessionCreator) createSessionDetails() *sessionDetails {
	return &sessionDetails{
		s.msg.SessionId,
		s.sessionToken,
		s.serverPubKey,
	}
}

func createNewSession() (*sessionDetails, error) {
	creator := &sessionCreator{
		pvtKey: readPemKey(),
	}

	err := creator.requestToken()
	if err != nil {
		return nil, err
	}

	err = creator.decryptSessionTokenWithPrivateKey()
	if err != nil {
		return nil, err
	}

	err = creator.decryptMessageWithSessionToken()
	if err != nil {
		return nil, err
	}

	err = creator.parseServerPublicKeyFromMessage()
	if err != nil {
		return nil, err
	}

	err = creator.verifyServerEncryptedSessionToken()
	if err != nil {
		return nil, err
	}

	return creator.createSessionDetails(), nil
}
