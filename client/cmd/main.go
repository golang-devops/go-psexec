package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-devops/go-psexec/client"
	"github.com/golang-devops/go-psexec/shared"
)

var (
	interactiveModeFlag = flag.Bool("interactive", false, "Interactive mode")
	serverFlag          = flag.String("server", "http://localhost:62677", "The endpoint server address")
	executorFlag        = flag.String("executor", "winshell", "The executor to use")
	clientPemFlag       = flag.String("client_pem", "client.pem", "The file path for the client pem (private+public) key file")
)

func handleRecovery() {
	if r := recover(); r != nil {
		log.Printf("ERROR: %s\n", getErrorStringFromRecovery(r))
	}
}

func execute(onFeedback func(fb string), server, executor, clientPemPath, exe string, args ...string) error {
	pvtKey, err := shared.ReadPemKey(clientPemPath)
	if err != nil {
		return fmt.Errorf("Cannot read client pem file, error: %s", err.Error())
	}

	c := client.New(pvtKey)

	session, err := c.RequestNewSession(server)
	if err != nil {
		return fmt.Errorf("Unable to create new session, error: %s", err.Error())
	}

	onFeedback(fmt.Sprintf("Using session id: %d\n", session.SessionId()))

	resp, err := session.StartExecRequest(&shared.ExecDto{executor, exe, args})
	if err != nil {
		return err
	}

	responseChannel, errChannel := resp.TextResponseChannel()

	allErrors := []error{}
outerFor:
	for {
		select {
		case feedbackLine, ok := <-responseChannel:
			if !ok {
				break outerFor
			}
			onFeedback(feedbackLine)
		case errLine, ok := <-errChannel:
			if !ok {
				break outerFor
			}
			allErrors = append(allErrors, errLine)
		}
	}

	if len(allErrors) > 0 {
		errStrs := []string{}
		for _, e := range allErrors {
			errStrs = append(errStrs, e.Error())
		}
		return fmt.Errorf("ERRORS WERE: %s", strings.Join(errStrs, "\\n"))
	}

	return nil
}

func main() {
	defer handleRecovery()

	flag.Parse()

	if *interactiveModeFlag {
		handleInteractiveMode()
		os.Exit(0)
	}

	exeAndArgs := flag.Args()
	if len(exeAndArgs) == 0 {
		log.Fatal("Need at least one additional argument")
	}

	var exe string
	var args []string = []string{}

	exe = exeAndArgs[0]
	if len(exeAndArgs) > 1 {
		args = exeAndArgs[1:]
	}

	onFeedback := func(fb string) {
		fmt.Println(fb)
	}
	err := execute(onFeedback, *serverFlag, *executorFlag, *clientPemFlag, exe, args...)
	if err != nil {
		log.Fatal(err)
	}
}
