package shared

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func ReadPemKey(filePath string) (*rsa.PrivateKey, error) {
	pem_data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pem_data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("No valid PEM data found")
	}

	pvtKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pvtKey.Precompute()

	err = pvtKey.Validate()
	if err != nil {
		return nil, err
	}

	return pvtKey, nil
}
