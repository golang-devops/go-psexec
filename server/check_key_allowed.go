package main

import (
	"crypto/rsa"
)

func (h *handler) checkPubKeyAllowed(pubKey *rsa.PublicKey) bool {
	for _, allowedPubKey := range h.AllowedPublicKeys {
		if eq, err := allowedPubKey.PublicKeyEquals(pubKey); eq {
			return true
		} else if err != nil {
			h.logger.Warningf("Failed to check public key allowed for name '%s', error: %s", allowedPubKey.Name, err.Error())
		}
	}
	return false
}
