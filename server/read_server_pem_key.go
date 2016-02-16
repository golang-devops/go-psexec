package main

import (
	"crypto/rsa"

	"github.com/golang-devops/go-psexec/shared"
)

func readPemKey() (*rsa.PrivateKey, error) {
	return shared.ReadPemKey(*serverPemFlag)
}
