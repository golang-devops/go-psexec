package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mozillazg/request"
	"log"
	"net/http"
	"strings"

	"github.com/golang-devops/go-psexec/shared"
)

//This should be the same as on the server...
const SigningKey = "somethingsupersecret"

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

	c := new(http.Client)

	jwtToken := generateToken()
	fmt.Println("Using JWT token: " + jwtToken)

	req := request.NewRequest(c)
	req.Headers["Authorization"] = "Bearer " + jwtToken

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

	req.Json = shared.Dto{*executorFlag, exe, args}

	url := strings.Trim(*serverFlag, "/") + "/auth/exec"

	resp, err := req.Post(url)
	checkError(err)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	defer resp.Body.Close()
}
