package shared

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func ReadPemKey(filePath string) *rsa.PrivateKey {
	pem_data, err := ioutil.ReadFile(filePath)
	checkError(err)

	block, _ := pem.Decode(pem_data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("No valid PEM data found")
	}

	pvtKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	checkError(err)

	pvtKey.Precompute()

	err = pvtKey.Validate()
	checkError(err)

	return pvtKey
}
