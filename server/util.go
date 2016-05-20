package main

import (
	"fmt"
	"net"

	"github.com/labstack/echo/engine"
)

func getErrorStringFromRecovery(r interface{}) string {
	errStr := ""
	switch t := r.(type) {
	case error:
		errStr = t.Error()
		break
	default:
		errStr = fmt.Sprintf("%#v", r)
	}
	return errStr
}

func getIPFromRequest(r engine.Request) string {
	if ipProxy := r.Header().Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddress())
	return ip
}

func getHostNamesFromIP(ip string) ([]string, error) {
	hostNames, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}
	return hostNames, nil
}
