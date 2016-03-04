package main

import (
	"fmt"
	"net"
	"net/http"
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

func getIPFromRequest(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getHostNamesFromIP(ip string) ([]string, error) {
	hostNames, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}
	return hostNames, nil
}
