package checksums

import (
	"encoding/hex"
	"fmt"
)

//ChecksumResult just contains the result of some checksum (md5) and gives us easy access to encoding it as hex or bytes, etc
type ChecksumResult struct{ b []byte }

//Bytes will return the checksum as a byte slice
func (c *ChecksumResult) Bytes() []byte {
	return c.b
}

//HexString will return the checksum encoded as hex
func (c *ChecksumResult) HexString() string {
	return hex.EncodeToString(c.b)
}

//NewChecksumResultFromHex creates a new ChecksumResult from the input hex string
func NewChecksumResultFromHex(hexStr string) (*ChecksumResult, error) {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode hex string '%s', error: %s", hexStr, err.Error())
	}
	return &ChecksumResult{decoded}, nil
}
