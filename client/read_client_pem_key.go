package main

import (
	"crypto/rsa"
	"os"
	"path/filepath"

	"github.com/golang-devops/go-psexec/shared"
)

func readPemKey() *rsa.PrivateKey {
	curExePath, err := filepath.Abs(os.Args[0])
	checkError(err)

	pemPath := filepath.Join(filepath.Dir(curExePath), "client.pem")

	return shared.ReadPemKey(pemPath)
}
