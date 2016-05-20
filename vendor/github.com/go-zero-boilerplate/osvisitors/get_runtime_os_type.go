package osvisitors

import "runtime"

func GetRuntimeOsType() (OsType, error) {
	return ParseFromName(runtime.GOOS)
}
