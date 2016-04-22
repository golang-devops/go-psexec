package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-devops/go-psexec/shared/tar_io"

	"github.com/go-zero-boilerplate/more_goconvey_assertions"

	"github.com/golang-devops/go-psexec/shared"
	. "github.com/smartystreets/goconvey/convey"
)

func testingGetNewSessionFileSystem() (SessionFileSystem, error) {
	clientPemFile := os.ExpandEnv(`$GOPATH/src/github.com/golang-devops/go-psexec/client/cmd/client.pem`)
	//TODO: Explicitly start a localhost server too and perhaps use another port to not conflict with "default" port of actual server. Refer to server `main_test.go` file method `TestHighLoad` in the `a.Run(logger)` call for an example
	serverUrl := "http://localhost:62677"

	clientPvtKey, err := shared.ReadPemKey(clientPemFile)
	if err != nil {
		return nil, err
	}

	client := New(clientPvtKey)

	session, err := client.RequestNewSession(serverUrl)
	if err != nil {
		return nil, err
	}

	return NewSessionFileSystem(session), nil
}

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

func TestUploadTarDirectory(t *testing.T) {
	Convey("Test Uploading a tar directory", t, func() {
		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		localAbsTestdataDir, err := filepath.Abs("testdata/upload-tar-dir")
		So(err, ShouldBeNil)

		localFileList, err := listFilesInDir(localAbsTestdataDir)
		So(err, ShouldBeNil)
		So(len(localFileList), ShouldBeGreaterThan, 0)

		// TODO: Make it more obvious that this directory is "populated" with the server which might be another machine
		// but for tests we plan to make the server startup locally on a different port so we can interact with it
		tempRemoteBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempRemoteBasePath)

		//The dir already exists due to the TempDir method
		for _, localFullFilePath := range localFileList {
			relPath := localFullFilePath[len(localAbsTestdataDir)+1:]
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relPath)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)
		}

		dirTarProvider := tar_io.Factories.TarProvider.Dir(localAbsTestdataDir, "")
		err = sessionFileSystem.UploadTar(dirTarProvider, tempRemoteBasePath)
		So(err, ShouldBeNil)

		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, true)
		for _, localFullFilePath := range localFileList {
			relPath := localFullFilePath[len(localAbsTestdataDir)+1:]
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relPath)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, true)
			So(checkFilePropertiesEqual(localFullFilePath, fullTempRemotePath), ShouldBeNil)
		}

		err = os.RemoveAll(tempRemoteBasePath)
		So(err, ShouldBeNil)
		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		for _, localFullFilePath := range localFileList {
			relPath := localFullFilePath[len(localAbsTestdataDir)+1:]
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relPath)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)
		}
	})
}

func TestUploadTarFile(t *testing.T) {
	Convey("Test Uploading a tar directory", t, func() {
		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		localAbsTestdataDir, err := filepath.Abs("testdata/upload-tar-file")
		So(err, ShouldBeNil)

		localFileList, err := listFilesInDir(localAbsTestdataDir)
		So(err, ShouldBeNil)
		So(len(localFileList), ShouldBeGreaterThan, 0)

		// TODO: Make it more obvious that this directory is "populated" with the server which might be another machine
		// but for tests we plan to make the server startup locally on a different port so we can interact with it
		tempRemoteBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempRemoteBasePath)

		for _, localFullFilePath := range localFileList {
			relPath := localFullFilePath[len(localAbsTestdataDir)+1:]
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relPath)

			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)

			fileTarProvider := tar_io.Factories.TarProvider.File(localFullFilePath)
			err = sessionFileSystem.UploadTar(fileTarProvider, fullTempRemotePath)
			So(err, ShouldBeNil)

			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, true)
			So(checkFilePropertiesEqual(localFullFilePath, fullTempRemotePath), ShouldBeNil)
		}

		err = os.RemoveAll(tempRemoteBasePath)
		So(err, ShouldBeNil)
		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		for _, localFullFilePath := range localFileList {
			relPath := localFullFilePath[len(localAbsTestdataDir)+1:]
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relPath)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)
		}
	})
}

func TestDownloadTarDirectory(t *testing.T) {
	Convey("Test Downloading a tar directory", t, func() {
		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		// TODO: Make it more obvious that this directory is "populated" with the server which might be another machine
		// but for tests we plan to make the server startup locally on a different port so we can interact with it
		remoteAbsTestdataDir, err := filepath.Abs("testdata/download-tar-dir")
		So(err, ShouldBeNil)

		remoteFileList, err := listFilesInDir(remoteAbsTestdataDir)
		So(err, ShouldBeNil)
		So(len(remoteFileList), ShouldBeGreaterThan, 0)

		tempLocalBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempLocalBasePath)

		//The dir already exists due to the TempDir method
		for _, remoteFullFilePath := range remoteFileList {
			relPath := remoteFullFilePath[len(remoteAbsTestdataDir)+1:]
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relPath)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)
		}

		tarReceiver := tar_io.Factories.TarReceiver.Dir(tempLocalBasePath)
		err = sessionFileSystem.DownloadTar(remoteAbsTestdataDir, nil, tarReceiver)
		So(err, ShouldBeNil)

		So(tempLocalBasePath, more_goconvey_assertions.AssertDirectoryExistance, true)
		for _, remoteFullFilePath := range remoteFileList {
			relPath := remoteFullFilePath[len(remoteAbsTestdataDir)+1:]
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relPath)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, true)
			So(checkFilePropertiesEqual(remoteFullFilePath, fullTempLocalPath), ShouldBeNil)
		}

		err = os.RemoveAll(tempLocalBasePath)
		So(err, ShouldBeNil)
		So(tempLocalBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		for _, remoteFullFilePath := range remoteFileList {
			relPath := remoteFullFilePath[len(remoteAbsTestdataDir)+1:]
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relPath)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)
		}
	})
}

func TestDownloadTarFile(t *testing.T) {
	Convey("Test Downloading a tar file", t, func() {
		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		// TODO: Make it more obvious that this directory is "populated" with the server which might be another machine
		// but for tests we plan to make the server startup locally on a different port so we can interact with it
		remoteAbsTestdataDir, err := filepath.Abs("testdata/download-tar-file")
		So(err, ShouldBeNil)

		remoteFileList, err := listFilesInDir(remoteAbsTestdataDir)
		So(err, ShouldBeNil)
		So(len(remoteFileList), ShouldBeGreaterThan, 0)

		tempLocalBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempLocalBasePath)

		//The dir already exists due to the TempDir method
		for _, remoteFullFilePath := range remoteFileList {
			relPath := remoteFullFilePath[len(remoteAbsTestdataDir)+1:]
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relPath)

			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)

			tarReceiver := tar_io.Factories.TarReceiver.File(fullTempLocalPath)
			err = sessionFileSystem.DownloadTar(remoteAbsTestdataDir, nil, tarReceiver)
			So(err, ShouldBeNil)

			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, true)
			So(checkFilePropertiesEqual(remoteFullFilePath, fullTempLocalPath), ShouldBeNil)
		}

		err = os.RemoveAll(tempLocalBasePath)
		So(err, ShouldBeNil)
		So(tempLocalBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		for _, remoteFullFilePath := range remoteFileList {
			relPath := remoteFullFilePath[len(remoteAbsTestdataDir)+1:]
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relPath)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)
		}
	})
}
