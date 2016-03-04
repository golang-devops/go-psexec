package service

import (
	"fmt"
	"os/user"
)

func getCurrentUserName() string {
	u, _ := user.Current()
	if u != nil {
		return u.Username
	}
	return ""
}

func getErrorStringFromRecovery(r interface{}) string {
	str := ""
	switch t := r.(type) {
	case error:
		str = t.Error()
		break
	default:
		str = fmt.Sprintf("%#v", r)
	}
	return str
}
