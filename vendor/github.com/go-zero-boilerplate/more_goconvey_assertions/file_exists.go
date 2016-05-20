package more_goconvey_assertions

import (
	"fmt"
	"reflect"

	"github.com/go-zero-boilerplate/path_utils"
)

var AssertFileExistance = assertFileExistance
var AssertDirectoryExistance = assertDirectoryExistance

func assertFileExistance(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	filePath, filePathIsString := actual.(string)
	mustExist, mustExistIsBool := expected[0].(bool)

	if !filePathIsString {
		return fmt.Sprintf(shouldBeString, reflect.TypeOf(actual))
	}

	if !mustExistIsBool {
		return fmt.Sprintf(shouldBeBool, reflect.TypeOf(expected[0]))
	}

	if fileExist, err := path_utils.FileExists(filePath); err != nil {
		return err.Error()
	} else if fileExist != mustExist {
		if mustExist {
			return fmt.Sprintf(shouldHaveBeenTrue, fileExist)
		} else {
			return fmt.Sprintf(shouldHaveBeenFalse, fileExist)
		}
	}

	return success
}

func assertDirectoryExistance(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	dirPath, dirPathIsString := actual.(string)
	mustExist, mustExistIsBool := expected[0].(bool)

	if !dirPathIsString {
		return fmt.Sprintf(shouldBeString, reflect.TypeOf(actual))
	}

	if !mustExistIsBool {
		return fmt.Sprintf(shouldBeBool, reflect.TypeOf(expected[0]))
	}

	if dirExist, err := path_utils.DirectoryExists(dirPath); err != nil {
		return err.Error()
	} else if dirExist != mustExist {
		if mustExist {
			return fmt.Sprintf(shouldHaveBeenTrue, dirExist)
		} else {
			return fmt.Sprintf(shouldHaveBeenFalse, dirExist)
		}
	}

	return success
}
