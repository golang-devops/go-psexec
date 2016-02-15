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

func (a *AllowedPublicKey) PublicKeyEquals(otherPubKey *rsa.PublicKey) bool {
	thisBytes, err := x509.MarshalPKIXPublicKey(a.PublicKey)
	checkError(err)

	otherBytes, err := x509.MarshalPKIXPublicKey(otherPubKey)
	checkError(err)

	return bytes.Equal(thisBytes, otherBytes)
}
