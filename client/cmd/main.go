package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-devops/go-psexec/client"
	"github.com/golang-devops/go-psexec/shared"
)

var (
	version = "0.0.2"
)

var (
	interactiveModeFlag = flag.Bool("interactive", false, "Interactive mode")
	fireAndForgetFlag   = flag.Bool("fire-and-forget", false, "Fire and forget (run the process in the background)")
	serverFlag          = flag.String("server", "http://localhost:62677", "The endpoint server address")
	executorFlag        = flag.String("executor", "winshell", "The executor to use")
	clientPemFlag       = flag.String("client_pem", "client.pem", "The file path for the client pem (private+public) key file")
	timeoutSeconds      = flag.Int64("timeout-seconds", 0, "The timeout to use when making calls")
)

func handleRecovery() {
	if r := recover(); r != nil {
		log.Printf("ERROR: %s\n", getErrorStringFromRecovery(r))
	}
}

func getExecutorExecRequestBuilder(session client.Session, executor string) client.SessionExecRequestBuilder {
	builder := session.ExecRequestBuilder()
	if executor == "winshell" {
		return builder.Winshell()
	} else if executor == "bash" {
		return builder.Bash()
	}
	panic("Unsupported client executor: '" + executor + "'")
}

func execute(fireAndForget bool, onFeedback func(fb string), server, executor, clientPemPath string, exe string, args ...string) error {
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

	workingDir := ""
	builder := getExecutorExecRequestBuilder(session, executor)
	resp, err := builder.Exe(exe).WorkingDir(workingDir).Args(args...).BuildAndDoRequest()
	if err != nil {
		return err
	}

	if fireAndForget {
		onFeedback("Using fire-and-forget mode so the command will continue in the background.")
		return nil
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

	fmt.Println("VERSION:", version)
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

	if *timeoutSeconds > 0 {
		go func() {
			time.Sleep(time.Duration(*timeoutSeconds) * time.Second)
			fmt.Printf("Timeout of %d seconds reached, forcefully aborting the client process\n", *timeoutSeconds)
			os.Exit(1)
		}()
	}

	onFeedback := func(fb string) {
		fmt.Println(fb)
	}
	err := execute(*fireAndForgetFlag, onFeedback, *serverFlag, *executorFlag, *clientPemFlag, exe, args...)
	if err != nil {
		log.Fatal(err)
	}
}
