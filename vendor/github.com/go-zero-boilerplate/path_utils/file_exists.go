package path_utils

import "os"

func FileExists(file string) (bool, error) {
	stats, err := os.Stat(file)
	if err == nil {
		if !stats.IsDir() {
			return true, nil
		} else {
			return false, nil
		}
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
