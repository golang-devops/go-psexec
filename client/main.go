package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mozillazg/request"
	"log"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
)

var (
	serverFlag   = flag.String("server", "http://localhost:62677", "The endpoint server address")
	executorFlag = flag.String("executor", "winshell", "The executor to use")
)

func handleRecovery() {
	if r := recover(); r != nil {
		log.Printf("ERROR: %s\n", getErrorStringFromRecovery(r))
	}
}

func main() {
	defer handleRecovery()

	flag.Parse()

	session, err := createNewSession()
	checkError(err)

	fmt.Printf("Using session id: %s\n", session.SessionId)

	exeAndArgs := flag.Args()
	if len(exeAndArgs) == 0 {
		panic("Need at least one additional argument")
	}

	var exe string
	var args []string = []string{}

	exe = exeAndArgs[0]
	if len(exeAndArgs) > 1 {
		args = exeAndArgs[1:]
	}

	c := new(http.Client)
	req := request.NewRequest(c)
	//req.Headers["Authorization"] = "Bearer " + token
	req.Json = shared.ExecDto{*executorFlag, exe, args}

	url := combineServerUrl("/auth/exec")

	resp, err := req.Post(url)
	checkError(err)

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
