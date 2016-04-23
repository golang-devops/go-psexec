package testing_utils

import (
	"fmt"
	"os"
)

func CheckFilePropertiesEqual(filePath1, filePath2 string) error {
	file1Info, err := os.Stat(filePath1)
	if err != nil {
		return fmt.Errorf("Cannot get file '%s' stats, error: %s", filePath1, err.Error())
	}
	file2Info, err := os.Stat(filePath2)
	if err != nil {
		return fmt.Errorf("Cannot get file '%s' stats, error: %s", filePath2, err.Error())
	}

	timestampFormat := "2006-01-02 15:04:05"
	timestamp1 := file1Info.ModTime().Format(timestampFormat)
	timestamp2 := file2Info.ModTime().Format(timestampFormat)
	if timestamp1 != timestamp2 {
		return fmt.Errorf("ModTime of file '%s' (%s) differs from file '%s' (%s)", filePath1, timestamp1, filePath2, timestamp2)
	}

	if file1Info.Size() != file2Info.Size() {
		return fmt.Errorf("Size of file '%s' (%d) differs from file '%s' (%d)", filePath1, file1Info.Size(), filePath2, file2Info.Size())
	}

	return nil
}
