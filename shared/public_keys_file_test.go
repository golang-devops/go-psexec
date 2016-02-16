package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadAllowedPublicKeysFile(t *testing.T) {
	Convey("Test loading of allowed public keys file", t, func() {
		allowedKeysPath, err := testsGetFilePath("allowed_keys")
		So(err, ShouldBeNil)

		allowedKeys, err := LoadAllowedPublicKeysFile(allowedKeysPath)
		So(err, ShouldBeNil)

		So(len(allowedKeys), ShouldEqual, 2)
		So(allowedKeys[0].Name, ShouldEqual, "Sender 1")
		So(allowedKeys[1].Name, ShouldEqual, "Another key")

		senderPvtKey, err := testsReadPemKey("sender.pem")
		So(err, ShouldBeNil)
		recipientPvtKey, err := testsReadPemKey("recipient.pem")
		So(err, ShouldBeNil)

		//The first allowed key equals the "sender pem key"
		equalSender, err := allowedKeys[0].PublicKeyEquals(&senderPvtKey.PublicKey)
		So(err, ShouldBeNil)
		So(equalSender, ShouldBeTrue)

		equalRecipient, err := allowedKeys[0].PublicKeyEquals(&recipientPvtKey.PublicKey)
		So(err, ShouldBeNil)
		So(equalRecipient, ShouldBeFalse)

		//The second allowed key is just a random key
		equalSender, err = allowedKeys[1].PublicKeyEquals(&senderPvtKey.PublicKey)
		So(err, ShouldBeNil)
		So(equalSender, ShouldBeFalse)

		equalRecipient, err = allowedKeys[1].PublicKeyEquals(&recipientPvtKey.PublicKey)
		So(err, ShouldBeNil)
		So(equalRecipient, ShouldBeFalse)
	})
}
