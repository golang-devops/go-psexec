package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/ayufan/golang-kardianos-service"
	"github.com/labstack/echo"
	"io"
	"net"
	"net/http"
)

type handler struct {
	logger     service.Logger
	privateKey *rsa.PrivateKey
}

func (h *handler) deserializeBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(dest)
}

func (h *handler) getPublicKeyBytes() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&h.privateKey.PublicKey)
}

func (h *handler) getAuthenticatedSessionToken(c *echo.Context) (*sessionToken, error) {
	sessionIdInterface := c.Get("session-id")
	sessionId, ok := sessionIdInterface.(int)
	if !ok {
		return nil, fmt.Errorf("Context session-id invalid format '%#v'", sessionIdInterface)
	}

	token, ok := tokenStore.GetSessionToken(sessionId)
	if !ok {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}

func (h *handler) getIPFromRequest(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (h *handler) getHostNamesFromIP(ip string) []string {
	hostNames, err := net.LookupAddr(ip)
	if err != nil {
		h.logger.Warningf("Unable to find hostname(s) for IP '%s', error: %s", ip, err.Error())
		return []string{}
	}
	return hostNames
}
