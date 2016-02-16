package shared

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	Convey("Generating a private key", t, func() {
		pvtKey, err := GeneratePrivateKey()
		So(err, ShouldBeNil)
		So(pvtKey, ShouldNotBeNil)
	})
}

func TestPrintPemFilePublicKeyAsHex(t *testing.T) {
	//No tests for now
}
