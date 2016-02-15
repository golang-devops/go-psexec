package main

import (
	"crypto/rsa"
)

func (h *handler) checkPubKeyAllowed(pubKey *rsa.PublicKey) bool {
	for _, allowedPubKey := range h.AllowedPublicKeys {
		if allowedPubKey.PublicKeyEquals(pubKey) {
			return true
		}
	}
	return false
}
