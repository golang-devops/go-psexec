package shared

import (
	"crypto/rsa"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func testsGetFilePath(fileName string) (string, error) {
	return filepath.Abs("./testdata/" + fileName)
}

func testsReadPemKey(keyName string) (*rsa.PrivateKey, error) {
	pemPath, err := testsGetFilePath(keyName)
	if err != nil {
		return nil, err
	}

	return ReadPemKey(pemPath)
}

func TestEncryptWithPublicKey(t *testing.T) {
	Convey("Testing encryption using public key", t, func() {
		recipientPvtKey, err := testsReadPemKey("recipient.pem")
		So(err, ShouldBeNil)
		recipientPubKey := &recipientPvtKey.PublicKey

		senderPvtKey, err := testsReadPemKey("sender.pem")
		So(err, ShouldBeNil)
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
