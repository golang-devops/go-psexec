package shared

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

var (
	//TODO: Not sure about always using this label
	label   = []byte("")
	newhash = crypto.SHA256
	opts    = &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto} // TODO: Is this `PSSSaltLengthAuto` fine?
)

func GenerateMessageSignature(senderPrivKey *rsa.PrivateKey, plainText []byte) ([]byte, error) {
	pssh := newhash.New()
	_, err := pssh.Write(plainText)
	if err != nil {
		return nil, err
	}

	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, senderPrivKey, newhash, hashed, opts)
	return signature, err
}

func EncryptWithPublicKey(recipientPubKey *rsa.PublicKey, plainText []byte) (cipher []byte, failure error) {
	sha_hash := sha256.New()

	encrypted, err := rsa.EncryptOAEP(sha_hash, rand.Reader, recipientPubKey, plainText, label)

	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func VerifySenderMessage(senderPublicKey *rsa.PublicKey, plainText, signature []byte) error {
	pssh := newhash.New()
	_, err := pssh.Write(plainText)
	if err != nil {
		return err
	}

	hashed := pssh.Sum(nil)

	return rsa.VerifyPSS(senderPublicKey, newhash, hashed, signature, opts)
}

func DecryptWithPrivateKey(recipientPvtKey *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	sha_hash := sha256.New()

	decrypted, err := rsa.DecryptOAEP(sha_hash, rand.Reader, recipientPvtKey, cipher, label)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
