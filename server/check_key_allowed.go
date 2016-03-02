package main

import (
	"crypto/rsa"
)

func (h *handler) checkPubKeyAllowed(pubKey *rsa.PublicKey) bool {
	isAllowed, warnings := allowedKeysStore.CheckAllowed(pubKey)
	for _, warn := range warnings {
		h.logger.Warning(warn)
	}
	return isAllowed
}
