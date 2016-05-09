package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/francoishill/afero"

	"github.com/go-zero-boilerplate/more_goconvey_assertions"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/golang-devops/go-psexec/services/encoding/checksums"
	"github.com/golang-devops/go-psexec/services/filepath_summary"
	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/tar_io"
	"github.com/golang-devops/go-psexec/shared/testing_utils"
)

func testingGetNewSessionFileSystem() (SessionFileSystem, error) {
	clientPemFile := os.ExpandEnv(`$GOPATH/src/github.com/golang-devops/go-psexec/client/testdata/test_client.pem`)
	//TODO: Explicitly start a localhost server too and perhaps use another port to not conflict with "default" port of actual server. Refer to server `main_test.go` file method `TestHighLoad` in the `a.Run(logger)` call for an example
	serverUrl := "http://localhost:62677"

	clientPvtKey, err := shared.ReadPemKey(clientPemFile)
	if err != nil {
		return nil, err
	}

	client := New(clientPvtKey)

	tmpSession, err := client.RequestNewSession(serverUrl)
	if err != nil {
		return nil, err
	}

	return NewSessionFileSystem(tmpSession.(*session)), nil
}

type copiedTestDataDetails struct {
	RemoteFS        *testing_utils.RemoteFileSystem
	BaseTempDir     string
	CopiedToTempDir string
}

func (c *copiedTestDataDetails) Cleanup() error {
	return c.RemoteFS.RemoveAll(c.BaseTempDir)
}

func copyTestDataToTempDir(testData *testing_utils.TestDataContainer, copyToRelSubDir string) (*copiedTestDataDetails, error) {
	//This method will copy the dir and files and confirm they exist
	remoteFS := testing_utils.NewTestingRemoteFileSystem()

	baseTempDir, err := remoteFS.TempDir()
	if err != nil {
		return nil, err
	}
	copiedToTempDir := filepath.Join(baseTempDir, copyToRelSubDir)
	if err = afero.CopyDir(testData.TestDataBaseFs, "", remoteFS, copiedToTempDir); err != nil {
		return nil, err
	}
	return &copiedTestDataDetails{RemoteFS: remoteFS, BaseTempDir: baseTempDir, CopiedToTempDir: copiedToTempDir}, nil
}

func TestUploadTarDirectory(t *testing.T) {
	Convey("Test Uploading a tar directory", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		tempRemoteBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempRemoteBasePath)

		//The dir already exists due to the TempDir method
		testData.ForeachRelativeFile(func(relFile string) {
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relFile)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)
		})

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		dirTarProvider := tar_io.Factories.TarProvider.Dir(testData.FullDir, "")
		err = sessionFileSystem.UploadTar(dirTarProvider, tempRemoteBasePath, true)
		So(err, ShouldBeNil)

		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, true)
		testData.ForeachRelativeFile(func(relFile string) {
			localFullFilePath := filepath.Join(testData.FullDir, relFile)
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relFile)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, true)
			So(testing_utils.CheckFilePropertiesEqual(localFullFilePath, fullTempRemotePath), ShouldBeNil)
		})

		err = os.RemoveAll(tempRemoteBasePath)
		So(err, ShouldBeNil)
		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		testData.ForeachRelativeFile(func(relFile string) {
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relFile)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)
		})
	})
}

func TestUploadTarFile(t *testing.T) {
	Convey("Test Uploading a tar directory", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		tempRemoteBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempRemoteBasePath)

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		testData.ForeachRelativeFile(func(relFile string) {
			localFullFilePath := filepath.Join(testData.FullDir, relFile)
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relFile)

			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)

			fileTarProvider := tar_io.Factories.TarProvider.File(localFullFilePath)
			err = sessionFileSystem.UploadTar(fileTarProvider, fullTempRemotePath, false)
			So(err, ShouldBeNil)

			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, true)
			So(testing_utils.CheckFilePropertiesEqual(localFullFilePath, fullTempRemotePath), ShouldBeNil)
		})

		err = os.RemoveAll(tempRemoteBasePath)
		So(err, ShouldBeNil)
		So(tempRemoteBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		testData.ForeachRelativeFile(func(relFile string) {
			fullTempRemotePath := filepath.Join(tempRemoteBasePath, relFile)
			So(fullTempRemotePath, more_goconvey_assertions.AssertFileExistance, false)
		})
	})
}

func TestDownloadTarDirectory(t *testing.T) {
	Convey("Test Downloading a tar directory", t, func() {
		testDataAsMockRemote, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		tempLocalBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempLocalBasePath)

		//The dir already exists due to the TempDir method
		testDataAsMockRemote.ForeachRelativeFile(func(relFile string) {
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relFile)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)
		})

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		tarReceiver := tar_io.Factories.TarReceiver.Dir(tempLocalBasePath)
		err = sessionFileSystem.DownloadTar(testDataAsMockRemote.FullDir, nil, tarReceiver)
		So(err, ShouldBeNil)

		So(tempLocalBasePath, more_goconvey_assertions.AssertDirectoryExistance, true)
		testDataAsMockRemote.ForeachRelativeFile(func(relFile string) {
			remoteFullFilePath := filepath.Join(testDataAsMockRemote.FullDir, relFile)
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relFile)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, true)
			So(testing_utils.CheckFilePropertiesEqual(remoteFullFilePath, fullTempLocalPath), ShouldBeNil)
		})

		err = os.RemoveAll(tempLocalBasePath)
		So(err, ShouldBeNil)
		So(tempLocalBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		testDataAsMockRemote.ForeachRelativeFile(func(relFile string) {
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relFile)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)
		})
	})
}

func TestDownloadTarFile(t *testing.T) {
	Convey("Test Downloading a tar file", t, func() {
		testDataAsMockRemote, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		tempLocalBasePath, err := ioutil.TempDir(os.TempDir(), "gopsexec-client-test-")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempLocalBasePath)

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		//The dir already exists due to the TempDir method
		testDataAsMockRemote.ForeachRelativeFile(func(relFile string) {
			remoteFullFilePath := filepath.Join(testDataAsMockRemote.FullDir, relFile)
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relFile)

			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)

			tarReceiver := tar_io.Factories.TarReceiver.File(fullTempLocalPath)
			err = sessionFileSystem.DownloadTar(remoteFullFilePath, nil, tarReceiver)
			So(err, ShouldBeNil)

			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, true)
			So(testing_utils.CheckFilePropertiesEqual(remoteFullFilePath, fullTempLocalPath), ShouldBeNil)
		})

		err = os.RemoveAll(tempLocalBasePath)
		So(err, ShouldBeNil)
		So(tempLocalBasePath, more_goconvey_assertions.AssertDirectoryExistance, false)
		testDataAsMockRemote.ForeachRelativeFile(func(relFile string) {
			fullTempLocalPath := filepath.Join(tempLocalBasePath, relFile)
			So(fullTempLocalPath, more_goconvey_assertions.AssertFileExistance, false)
		})
	})
}

func TestDeleteDir(t *testing.T) {
	Convey("Test deletion of directory", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		copiedDetails, err := copyTestDataToTempDir(testData, "")
		So(err, ShouldBeNil)
		defer copiedDetails.Cleanup()

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		So(copiedDetails.BaseTempDir, more_goconvey_assertions.AssertDirectoryExistance, true)
		err = sessionFileSystem.Delete(copiedDetails.BaseTempDir)
		So(err, ShouldBeNil)
		So(copiedDetails.BaseTempDir, more_goconvey_assertions.AssertDirectoryExistance, false)

		testData.ForeachRelativeFile(func(relFile string) {
			fullRemoteFilePath := filepath.Join(copiedDetails.BaseTempDir, relFile)
			So(fullRemoteFilePath, more_goconvey_assertions.AssertFileExistance, false)
		})
	})
}

func TestDeleteFile(t *testing.T) {
	Convey("Test deletion of file", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		copiedDetails, err := copyTestDataToTempDir(testData, "")
		So(err, ShouldBeNil)
		defer copiedDetails.Cleanup()

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		testData.ForeachRelativeFile(func(relFile string) {
			fullRemoteFilePath := filepath.Join(copiedDetails.BaseTempDir, relFile)
			So(fullRemoteFilePath, more_goconvey_assertions.AssertFileExistance, true)
			err = sessionFileSystem.Delete(fullRemoteFilePath)
			So(err, ShouldBeNil)
			So(fullRemoteFilePath, more_goconvey_assertions.AssertFileExistance, false)
		})
	})
}

func TestMoveDir(t *testing.T) {
	Convey("Test deletion of directory", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		copiedDetails, err := copyTestDataToTempDir(testData, "orig")
		So(err, ShouldBeNil)
		defer copiedDetails.Cleanup()

		remoteOrigDir := copiedDetails.CopiedToTempDir
		remoteNewDir := filepath.Join(copiedDetails.BaseTempDir, "new")

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		So(remoteOrigDir, more_goconvey_assertions.AssertDirectoryExistance, true)
		testData.ForeachRelativeFile(func(relFile string) {
			fullOldRemoteFilePath := filepath.Join(remoteOrigDir, relFile)
			So(fullOldRemoteFilePath, more_goconvey_assertions.AssertFileExistance, true)
			fullNewRemoteFilePath := filepath.Join(remoteNewDir, relFile)
			So(fullNewRemoteFilePath, more_goconvey_assertions.AssertFileExistance, false)
		})
		So(remoteNewDir, more_goconvey_assertions.AssertDirectoryExistance, false)

		err = sessionFileSystem.Move(remoteOrigDir, remoteNewDir)
		So(err, ShouldBeNil)
		So(remoteOrigDir, more_goconvey_assertions.AssertDirectoryExistance, false)
		So(remoteNewDir, more_goconvey_assertions.AssertDirectoryExistance, true)

		testData.ForeachRelativeFile(func(relFile string) {
			fullOldRemoteFilePath := filepath.Join(remoteOrigDir, relFile)
			So(fullOldRemoteFilePath, more_goconvey_assertions.AssertFileExistance, false)
			fullNewRemoteFilePath := filepath.Join(remoteNewDir, relFile)
			So(fullNewRemoteFilePath, more_goconvey_assertions.AssertFileExistance, true)
		})
	})
}

func TestMoveFile(t *testing.T) {
	Convey("Test deletion of file", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		copiedDetails, err := copyTestDataToTempDir(testData, "orig")
		So(err, ShouldBeNil)
		defer copiedDetails.Cleanup()

		remoteOrigDir := copiedDetails.CopiedToTempDir
		remoteNewDir := filepath.Join(copiedDetails.BaseTempDir, "new")

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		So(remoteOrigDir, more_goconvey_assertions.AssertDirectoryExistance, true)
		testData.ForeachRelativeFile(func(relFile string) {
			fullOldRemoteFilePath := filepath.Join(remoteOrigDir, relFile)
			So(fullOldRemoteFilePath, more_goconvey_assertions.AssertFileExistance, true)
			fullNewRemoteFilePath := filepath.Join(remoteNewDir, relFile)
			So(fullNewRemoteFilePath, more_goconvey_assertions.AssertFileExistance, false)
		})
		So(remoteNewDir, more_goconvey_assertions.AssertDirectoryExistance, false)

		testData.ForeachRelativeFile(func(relFile string) {
			fullOldRemoteFilePath := filepath.Join(remoteOrigDir, relFile)
			So(fullOldRemoteFilePath, more_goconvey_assertions.AssertFileExistance, true)
			fullNewRemoteFilePath := filepath.Join(remoteNewDir, relFile)
			So(fullNewRemoteFilePath, more_goconvey_assertions.AssertFileExistance, false)
			err = sessionFileSystem.Move(fullOldRemoteFilePath, fullNewRemoteFilePath)
			So(err, ShouldBeNil)
		})

		//The old dir must still exist, we only moved the files
		So(remoteOrigDir, more_goconvey_assertions.AssertDirectoryExistance, true)
		So(remoteNewDir, more_goconvey_assertions.AssertDirectoryExistance, true)

		testData.ForeachRelativeFile(func(relFile string) {
			fullOldRemoteFilePath := filepath.Join(remoteOrigDir, relFile)
			So(fullOldRemoteFilePath, more_goconvey_assertions.AssertFileExistance, false)
			parentOldDir := filepath.Dir(fullOldRemoteFilePath)
			So(parentOldDir, more_goconvey_assertions.AssertDirectoryExistance, true)
			fullNewRemoteFilePath := filepath.Join(remoteNewDir, relFile)
			So(fullNewRemoteFilePath, more_goconvey_assertions.AssertFileExistance, true)
		})
	})
}

func TestStats(t *testing.T) {
	Convey("Test stats of remote file system", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		copiedDetails, err := copyTestDataToTempDir(testData, "")
		So(err, ShouldBeNil)
		defer copiedDetails.Cleanup()

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)

		testData.ForeachRelativePath(func(relPath string) {
			fullRemotePath := filepath.Join(copiedDetails.BaseTempDir, relPath)
			remoteStats, err := sessionFileSystem.Stats(fullRemotePath)
			So(err, ShouldBeNil)
			So(remoteStats, ShouldNotBeNil)

			localStats, err := copiedDetails.RemoteFS.Stat(filepath.Join(copiedDetails.BaseTempDir, relPath))
			So(err, ShouldBeNil)

			So(remoteStats.IsDir, ShouldEqual, localStats.IsDir())
			if !remoteStats.ModTime.Equal(localStats.ModTime()) {
				So(fmt.Errorf("Unexpected ModTime stamps remote '%s' vs local '%s'", remoteStats.ModTime.String(), localStats.ModTime().String()), ShouldBeNil)
			}
			So(remoteStats.Mode, ShouldEqual, localStats.Mode())
			So(remoteStats.Size, ShouldEqual, localStats.Size())
		})
	})
}

func getFileSummaryForFile(checksumSvc checksums.Service, filePath string) (*filepath_summary.FileSummary, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	checksumResult, err := checksumSvc.FileChecksum(filePath)
	if err != nil {
		return nil, err
	}
	return filepath_summary.NewFileSummary(filePath, fileInfo.ModTime(), checksumResult), nil
}

func TestDirSummary(t *testing.T) {
	Convey("Test dir-summary of remote file system", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		copiedDetails, err := copyTestDataToTempDir(testData, "")
		So(err, ShouldBeNil)
		defer copiedDetails.Cleanup()

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)
		checksumSvc := checksums.New()

		remoteOrigDir := copiedDetails.CopiedToTempDir
		summary, err := sessionFileSystem.DirSummary(remoteOrigDir)
		So(err, ShouldBeNil)
		So(summary, ShouldNotBeNil)

		So(summary.FlattenedFileSummaries, ShouldNotResemble, []*filepath_summary.FileSummary{})
		So(len(summary.FlattenedFileSummaries), ShouldNotEqual, 0)

		testData.ForeachRelativeFile(func(relPath string) {
			var actualFileSummary *filepath_summary.FileSummary
			for _, fileSummary := range summary.FlattenedFileSummaries {
				summaryRelPath := strings.TrimLeft(fileSummary.FullPath[len(remoteOrigDir):], "\\/")
				if summaryRelPath == relPath {
					actualFileSummary = fileSummary
					break
				}
			}
			if actualFileSummary == nil {
				So(fmt.Errorf("Did not find summary for file '%s'", relPath), ShouldBeNil)
			}

			fullRemoteFilePath := filepath.Join(copiedDetails.BaseTempDir, relPath)
			expectedSummary, err := getFileSummaryForFile(checksumSvc, fullRemoteFilePath)
			So(err, ShouldBeNil)

			So(actualFileSummary.Checksum.HexString(), ShouldEqual, expectedSummary.Checksum.HexString())
			So(actualFileSummary.ModTime.Equal(expectedSummary.ModTime), ShouldBeTrue)
		})

		So(len(summary.FlattenedFileSummaries), ShouldEqual, len(testData.RelativeFiles))
	})
}

func TestFileSummary(t *testing.T) {
	Convey("Test file-summary of remote file system", t, func() {
		testData, err := testing_utils.NewTestDataContainer("testdata/copy-files", 6)
		So(err, ShouldBeNil)

		copiedDetails, err := copyTestDataToTempDir(testData, "")
		So(err, ShouldBeNil)
		defer copiedDetails.Cleanup()

		sessionFileSystem, err := testingGetNewSessionFileSystem()
		So(err, ShouldBeNil)
		checksumSvc := checksums.New()

		remoteOrigDir := copiedDetails.CopiedToTempDir
		testData.ForeachRelativeFile(func(relPath string) {
			actualFileSummary, err := sessionFileSystem.FileSummary(filepath.Join(remoteOrigDir, relPath))
			So(err, ShouldBeNil)
			So(actualFileSummary, ShouldNotBeNil)

			fullRemoteFilePath := filepath.Join(copiedDetails.CopiedToTempDir, relPath)
			expectedSummary, err := getFileSummaryForFile(checksumSvc, fullRemoteFilePath)
			So(err, ShouldBeNil)

			So(actualFileSummary.Checksum.HexString(), ShouldEqual, expectedSummary.Checksum.HexString())
			So(actualFileSummary.ModTime.Equal(expectedSummary.ModTime), ShouldBeTrue)
		})
	})
}
