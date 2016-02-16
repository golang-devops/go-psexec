package shared

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateKeyPairPemFile(outputPemFilePath string) error {
	pvtKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	pvtKey.Precompute()

	err = pvtKey.Validate()
	if err != nil {
		return err
	}

	pvtFile, err := os.OpenFile(outputPemFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer pvtFile.Close()

	return pem.Encode(pvtFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pvtKey),
	})
}

func GeneratePublicKeyFromPemFlag(inputPemFile string) error {
	privKey, err := ReadPemKey(inputPemFile)
	if err != nil {
		return err
	}

	pubPKIXBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return err
	}

	hexString := hex.EncodeToString(pubPKIXBytes)
	fmt.Println(hexString)
	return nil
}
