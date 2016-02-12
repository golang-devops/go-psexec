package shared

import (
	"crypto/rsa"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func testsReadPemKey(keyName string) *rsa.PrivateKey {
	pemPath, err := filepath.Abs("./testdata/" + keyName)
	checkError(err)

	return ReadPemKey(pemPath)
}

func TestEncryptWithPublicKey(t *testing.T) {
	Convey("Testing encryption using public key", t, func() {
		recipientPvtKey := testsReadPemKey("recipient.pem")
		recipientPubKey := &recipientPvtKey.PublicKey
		senderPvtKey := testsReadPemKey("sender.pem")
		senderPublicKey := &senderPvtKey.PublicKey
		message := []byte{122, 182, 172, 135, 64, 225, 230, 251, 55, 13, 230, 6, 164, 86, 124, 7, 180, 233, 94, 188, 172, 253, 142, 131, 166, 36, 43, 135, 70, 23, 128, 2}

		cipher, err := EncryptWithPublicKey(recipientPubKey, message)
		So(err, ShouldBeNil)

		signature, err := GenerateMessageSignature(senderPvtKey, message)
		So(err, ShouldBeNil)

		origMsg, err := DecryptWithPrivateKey(recipientPvtKey, cipher)
		So(err, ShouldBeNil)
		So(origMsg, ShouldResemble, message)

		err = VerifySenderMessage(senderPublicKey, origMsg, signature)
		So(err, ShouldBeNil)
	})
}
