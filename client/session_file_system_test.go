package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-zero-boilerplate/more_goconvey_assertions"

	"github.com/golang-devops/go-psexec/shared"
	. "github.com/smartystreets/goconvey/convey"
)

func listFilesInDir(dir string) ([]string, error) {
	files := []string{}
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, errParam error) error {
		if errParam != nil {
			return errParam
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	return files, nil
}

func checkFilePropertiesEqual(filePath1, filePath2 string) error {
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

func TestUploadDir(t *testing.T) {
	Convey("Test Uploading a directory", t, func() {
		clientPemFile := os.ExpandEnv(`$GOPATH/src/github.com/golang-devops/go-psexec/client/cmd/client.pem`)
		//TODO: Explicitly start a localhost server too and perhaps use another port to not conflict with "default" port of actual server. Refer to server `main_test.go` file method `TestHighLoad` in the `a.Run(logger)` call for an example
		serverUrl := "http://localhost:62677"

		clientPvtKey, err := shared.ReadPemKey(clientPemFile)
		So(err, ShouldBeNil)

		client := New(clientPvtKey)

		session, err := client.RequestNewSession(serverUrl)
		So(err, ShouldBeNil)

		sessionFileSystem := NewSessionFileSystem(session)

		localAbsTestdataDir, err := filepath.Abs("testdata/dir-reader")
		So(err, ShouldBeNil)

		localFileList, err := listFilesInDir(localAbsTestdataDir)
		So(err, ShouldBeNil)
		So(len(localFileList), ShouldBeGreaterThan, 0)

		// Make it more obvious that this directory is "populated" with the server which might be another machine
		// but for tests we plan to make the server startup locally on a different port so we can interact with it
		tempRemoteBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempRemoteBasePath)

		localTarReader := NewDirTarReader(localAbsTestdataDir, "", tempRemoteBasePath)
		err = sessionFileSystem.Upload(localTarReader)
		So(err, ShouldBeNil)

		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, true)
		for _, localFullFilePath := range localFileList {
			relPath := localFullFilePath[len(localAbsTestdataDir)+1:]
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relPath)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, true)
			So(checkFilePropertiesEqual(localFullFilePath, fullTempRemotePath), ShouldBeNil)
		}

		os.RemoveAll(tempRemoteBasePath)
		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
	})
}
