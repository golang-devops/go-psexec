package main

import (
	"crypto/rand"
	"crypto/rsa"
	"sync"

	"github.com/golang-devops/go-psexec/shared"
)

var (
	currentTokenId int
	tokenStore     = &TokenStore{tokens: make(map[int]*sessionToken)}
)

type TokenStore struct {
	sync.RWMutex
	tokens map[int]*sessionToken
}

type sessionToken struct {
	Token           []byte
	ClientPublicKey *rsa.PublicKey
}

func (s *sessionToken) DecryptWithSessionToken(cipher []byte) ([]byte, error) {
	return shared.DecryptSymmetric(s.Token, cipher)
}

/*func (s *sessionToken) NewEncryptedWriter(writer io.Writer) *shared.EncryptedWriterProxy {
	return shared.NewEncryptedWriterProxy(writer, s.Token)
}*/

func (t *TokenStore) NewSessionToken(clientPublicKey *rsa.PublicKey) (int, []byte, error) {
	t.Lock()
	defer t.Unlock()

	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return 0, nil, err
	}

	currentTokenId++

	t.tokens[currentTokenId] = &sessionToken{key, clientPublicKey}

	// The key length can be 32, 24, 16  bytes (OR in bits: 128, 192 or 256)
	return currentTokenId, key, nil
}

func (t *TokenStore) GetSessionToken(sessionId int) (*sessionToken, bool) {
	t.Lock()
	defer t.Unlock()

	tok, ok := t.tokens[sessionId]
	return tok, ok
}

func init() {
	currentTokenId = 1
}
