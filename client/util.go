package main

import (
	"fmt"
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
