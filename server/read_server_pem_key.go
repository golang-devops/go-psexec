package main

import (
	"crypto/rsa"
	"os"
	"path/filepath"

	"github.com/golang-devops/go-psexec/shared"
)

func readPemKey() (*rsa.PrivateKey, error) {
	curExePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil, err
	}

	pemPath := filepath.Join(filepath.Dir(curExePath), "server.pem")

	return shared.ReadPemKey(pemPath)
}
