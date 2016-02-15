package shared

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"io/ioutil"
	"strings"
)

func LoadAllowedPublicKeysFile(file string) (allowedKeys []*AllowedPublicKey) {
	fileBytes, err := ioutil.ReadFile(file)
	checkError(err)

	lines := strings.Split(string(fileBytes), "\n")

	for _, l := range lines {
		trimmedLine := strings.TrimSpace(l)
		if trimmedLine == "" {
			continue
		}

		split := strings.Split(trimmedLine, ":")

		name := strings.TrimSpace(split[0])
		hexKey := strings.TrimSpace(split[1])

		pubPKIXBytes, err := hex.DecodeString(hexKey)
		checkError(err)

		pubKeyInterface, err := x509.ParsePKIXPublicKey(pubPKIXBytes)
		checkError(err)

		rsaPublicKey, ok := pubKeyInterface.(*rsa.PublicKey)
		if !ok {
			panic("The server public-key received is in an incorrect format")
		}

		allowedKeys = append(allowedKeys,
			&AllowedPublicKey{
				name,
				rsaPublicKey,
			})
	}

	return
}
