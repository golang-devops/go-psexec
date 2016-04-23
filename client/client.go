package client

import (
	"crypto/rsa"
)

type Client struct {
	pvtKey *rsa.PrivateKey
}

func (c *Client) RequestNewSession(baseServerUrl string) (Session, error) {
	creator := &sessionCreator{
		pvtKey:        c.pvtKey,
		baseServerUrl: baseServerUrl,
	}

	err := creator.RequestToken()
	if err != nil {
		return nil, err
	}

	err = creator.DecryptSessionTokenWithPrivateKey()
	if err != nil {
		return nil, err
	}

	err = creator.DecryptMessageWithSessionToken()
	if err != nil {
		return nil, err
	}

	err = creator.ParseServerPublicKeyFromMessage()
	if err != nil {
		return nil, err
	}

	err = creator.VerifyServerEncryptedSessionToken()
	if err != nil {
		return nil, err
	}

	return creator.Create(), nil
}

func New(clientPrivateKey *rsa.PrivateKey) *Client {
	return &Client{
		clientPrivateKey,
	}
}
