package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-devops/go-psexec/shared"
	"sync"
)

var allowedKeysStore *AllowedKeysStore = &AllowedKeysStore{}

type AllowedKeysStore struct {
	sync.RWMutex
	keys []*shared.AllowedPublicKey
}

func (a *AllowedKeysStore) Clear() {
	a.Lock()
	defer a.Unlock()
	a.keys = []*shared.AllowedPublicKey{}
}

func (a *AllowedKeysStore) AppendKey(k *shared.AllowedPublicKey) {
	a.Lock()
	defer a.Unlock()
	a.keys = append(a.keys, k)
}

func (a *AllowedKeysStore) CheckAllowed(pubKey *rsa.PublicKey) (isAllowed bool, warnings []string) {
	a.Lock()
	defer a.Unlock()

	for _, allowedPubKey := range a.keys {
		if eq, err := allowedPubKey.PublicKeyEquals(pubKey); eq {
			return true, warnings
		} else if err != nil {
			warnings = append(warnings, fmt.Sprintf("Failed to check public key allowed for name '%s', error: %s", allowedPubKey.Name, err.Error()))
		}
	}

	return false, warnings
}
