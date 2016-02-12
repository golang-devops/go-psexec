package main

import (
	"crypto/rand"
	"crypto/rsa"
	"sync"

	"github.com/golang-devops/go-psexec/shared"
)

var (
	newTokenLock   *sync.RWMutex
	currentTokenId int
	tmpTokens      map[int]*sessionToken = map[int]*sessionToken{}
)

type sessionToken struct {
	Token           []byte
	ClientPublicKey *rsa.PublicKey
}

func (s *sessionToken) DecryptWithServerPrivateKey(serverPrivateKey *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	return shared.DecryptWithPrivateKey(serverPrivateKey, cipher)
}

func newSessionToken(clientPublicKey *rsa.PublicKey) (int, []byte, error) {
	newTokenLock.Lock()
	defer newTokenLock.Unlock()

	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return 0, nil, err
	}

	currentTokenId++

	tmpTokens[currentTokenId] = &sessionToken{key, clientPublicKey}

	// The key length can be 32, 24, 16  bytes (OR in bits: 128, 192 or 256)
	return currentTokenId, key, nil
}

func init() {
	newTokenLock = &sync.RWMutex{}
	currentTokenId = 1
}
