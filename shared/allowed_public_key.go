package shared

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
)

type AllowedPublicKey struct {
	Name      string
	PublicKey *rsa.PublicKey
}

func (a *AllowedPublicKey) PublicKeyEquals(otherPubKey *rsa.PublicKey) (bool, error) {
	thisBytes, err := x509.MarshalPKIXPublicKey(a.PublicKey)
	if err != nil {
		return false, err
	}

	otherBytes, err := x509.MarshalPKIXPublicKey(otherPubKey)
	if err != nil {
		return false, err
	}

	return bytes.Equal(thisBytes, otherBytes), nil
}
