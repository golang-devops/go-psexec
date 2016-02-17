package shared

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"
)

func LoadAllowedPublicKeysFile(file string) (allowedKeys []*AllowedPublicKey, returnErr error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

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
		if err != nil {
			return nil, err
		}

		pubKeyInterface, err := x509.ParsePKIXPublicKey(pubPKIXBytes)
		if err != nil {
			return nil, err
		}

		rsaPublicKey, ok := pubKeyInterface.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("The server public-key received is in an incorrect format")
		}

		allowedKeys = append(allowedKeys,
			&AllowedPublicKey{
				name,
				rsaPublicKey,
			})
	}

	returnErr = nil
	return
}
