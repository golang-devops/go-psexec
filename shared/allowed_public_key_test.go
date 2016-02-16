package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPublicKeyEquals(t *testing.T) {
	Convey("Testing the PublicKeyEquals method", t, func() {
		pvtKey1, err := testsReadPemKey("recipient.pem")
		So(err, ShouldBeNil)
		pvtKey2, err := testsReadPemKey("sender.pem")
		So(err, ShouldBeNil)

		allowedKey1 := &AllowedPublicKey{"Recipient", &pvtKey1.PublicKey}

		equalsKey2, err := allowedKey1.PublicKeyEquals(&pvtKey2.PublicKey)
		So(err, ShouldBeNil)
		So(equalsKey2, ShouldBeFalse)

		equalsKey1, err := allowedKey1.PublicKeyEquals(&pvtKey1.PublicKey)
		So(err, ShouldBeNil)
		So(equalsKey1, ShouldBeTrue)
	})
}
