package shared

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
)

//TODO: Not sure about always using this label
var label = []byte("")

func EncryptWithPublicKey(pubKey *rsa.PublicKey, message []byte) ([]byte, error) {
	md5_hash := md5.New()

	encrypted, err := rsa.EncryptOAEP(md5_hash, rand.Reader, pubKey, message, label)

	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func DecryptWithPrivateKey(pvtKey *rsa.PrivateKey, encryptedMsg []byte) ([]byte, error) {
	md5_hash := md5.New()
	decrypted, err := rsa.DecryptOAEP(md5_hash, rand.Reader, pvtKey, encryptedMsg, label)

	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
