package main

import (
	"bufio"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-devops/go-psexec/client"
	"github.com/golang-devops/go-psexec/shared"
)

func testsGetFilePath(fileName string) (string, error) {
	return filepath.Abs("../shared/testdata/" + fileName)
}

type tmpTestsLogger struct {
	sync.RWMutex
	ErrorList []string
}

func (t *tmpTestsLogger) Info(v ...interface{}) error {
	// fmt.Println(v...)
	return nil
}
func (t *tmpTestsLogger) Infof(frmt string, a ...interface{}) error {
	// fmt.Println(fmt.Sprintf(frmt, a...))
	return nil
}
func (t *tmpTestsLogger) Warning(v ...interface{}) error {
	t.Lock()
	defer t.Unlock()
	// fmt.Println(v...)
	t.ErrorList = append(t.ErrorList, fmt.Sprintln(v...))
	return nil
}
func (t *tmpTestsLogger) Warningf(frmt string, a ...interface{}) error {
	t.Warning(fmt.Sprintln(fmt.Sprintf(frmt, a...)))
	return nil
}
func (t *tmpTestsLogger) Error(v ...interface{}) error {
	t.Lock()
	defer t.Unlock()

	// fmt.Println(v...)
	t.ErrorList = append(t.ErrorList, fmt.Sprintln(v...))
	return nil
}
func (t *tmpTestsLogger) Errorf(frmt string, a ...interface{}) error {
	t.Error(fmt.Sprintln(fmt.Sprintf(frmt, a...)))
	return nil
}

func setupClient(clientPemFile string) (*client.Client, error) {
	pvtKey, err := shared.ReadPemKey(clientPemFile)
	if err != nil {
		return nil, err
	}
	return client.New(pvtKey), nil
}

func cleanFeedbackLine(line string) string {
	return strings.Trim(line, " \"'")
}

func doRequest(wg *sync.WaitGroup, logger *tmpTestsLogger, index int, cl *client.Client, serverBaseUrl string) {
	defer wg.Done()

	session, err := cl.RequestNewSession(serverBaseUrl)
	if err != nil {
		logger.Errorf("Index %d (RequestNewSession) err: %s", index, err.Error())
		return
	}

	echoStr := fmt.Sprintf("Hallo (%d)", index)
	resp, err := session.StartExecWinshellRequest("echo", echoStr)
	if err != nil {
		logger.Errorf("Index %d (StartExecWinshellRequest) err: %s", index, err.Error())
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	expectedFeedback := []string{echoStr, shared.RESPONSE_EOF}
	if len(lines) != len(expectedFeedback) {
		logger.Errorf("Index %d expected was %#v, but actual was %#v", index, expectedFeedback, lines)
		return
	}

	for i, expLine := range expectedFeedback {
		if cleanFeedbackLine(lines[i]) != cleanFeedbackLine(expLine) {
			logger.Errorf("Index %d expected was %#v, but actual was %#v", index, expectedFeedback, lines)
			return
		}
	}
}

func TestHighLoad(t *testing.T) {
	Convey("Test HighLoad", t, func() {
		logger := &tmpTestsLogger{ErrorList: []string{}}
		a := &app{}

		port := "64040"
		serverAddress := "localhost:" + port
		serverBaseUrl := "http://localhost:" + port

		serverPemPath, err := testsGetFilePath("recipient.pem")
		So(err, ShouldBeNil)
		allowedKeysPath, err := testsGetFilePath("allowed_keys")
		So(err, ShouldBeNil)
		clientPemPath, err := testsGetFilePath("sender.pem")
		So(err, ShouldBeNil)

		addressFlag = &serverAddress
		serverPemFlag = &serverPemPath
		allowedPublicKeysFileFlag = &allowedKeysPath
		go a.Run(logger)

		time.Sleep(500 * time.Millisecond) //Give server time to start
		cl, err := setupClient(clientPemPath)
		So(err, ShouldBeNil)

		num := 300
		var wg sync.WaitGroup
		wg.Add(num)
		for i := 0; i < num; i++ {
			go doRequest(&wg, logger, i, cl, serverBaseUrl)
		}
		wg.Wait()

		/*time.Sleep(1 * time.Second)
		a.srv.Stop(1 * time.Millisecond)*/

		for i, e := range logger.ErrorList {
			t.Errorf("ErrorList[%d]: %s", i, e)
		}
		So(len(logger.ErrorList), ShouldEqual, 0)
	})
}
