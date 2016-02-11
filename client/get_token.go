package main

import (
	"crypto/x509"
	"github.com/mozillazg/request"
	"io/ioutil"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
)

func getToken() (string, error) {
	pvtKey := readPemKey()

	pubPKIXBytes, err := x509.MarshalPKIXPublicKey(&pvtKey.PublicKey)
	if err != nil {
		return "", err
	}

	c := new(http.Client)
	req := request.NewRequest(c)
	req.Json = &shared.GetTokenRequestDto{pubPKIXBytes}

	url := combineServerUrl("/token")
	resp, err := req.Post(url)
	checkError(err)
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}
