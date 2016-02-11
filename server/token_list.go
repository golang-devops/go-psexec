package main

import (
	"errors"
	"fmt"
	"sync"
)

var (
	newTokenLock *sync.RWMutex
	tmpCounter   int
)

var tmpTokens map[string]string = map[string]string{
	"my-tok": "pub-key",
}

func getClientPubkeyFromToken(token string) (string, error) {
	if pubKey, ok := tmpTokens[token]; ok {
		return pubKey, nil
	}
	return "", errors.New("Invalid token")
}

func newSessionToken() string {
	//TODO: Implement better token
	newTokenLock.Lock()
	defer newTokenLock.Unlock()
	tmpCounter++
	return fmt.Sprintf("%d", tmpCounter)
}

func init() {
	newTokenLock = &sync.RWMutex{}
	tmpCounter = 1
}
