package shared

import (
	"crypto/rand"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEncryptWithSessionToken(t *testing.T) {
	Convey("Encryption with session token", t, func() {
		key := make([]byte, 32)

		_, err := rand.Read(key)
		So(err, ShouldBeNil)

		plaintext := []byte("Hallo there\nI am a message to be encrypted. Lets add some characters $#%^$^$&#$(*7,';dfsl[powqeriuy928354 3s21d6fa5s4987321234 s?><:}{P~!@#$%^&*()_+}")
		ciphertext, err := EncryptSymmetric(key, plaintext)
		So(err, ShouldBeNil)

		decrypted, err := DecryptSymmetric(key, ciphertext)
		So(err, ShouldBeNil)

		So(plaintext, ShouldResemble, decrypted)
	})
}
