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

func GenerateKeyPairPemFile(outputPemFilePath string) {
	pvtKey, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)

	pvtKey.Precompute()

	err = pvtKey.Validate()
	checkError(err)

	pvtFile, err := os.OpenFile(outputPemFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	checkError(err)
	defer pvtFile.Close()

	err = pem.Encode(pvtFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pvtKey),
	})
	checkError(err)
}

func GeneratePublicKeyFromPemFlag(inputPemFile string) {
	privKey := ReadPemKey(inputPemFile)

	pubPKIXBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	checkError(err)

	hexString := hex.EncodeToString(pubPKIXBytes)
	fmt.Println(hexString)
}
