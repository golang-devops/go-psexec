package main

import (
	"flag"
	"github.com/zero-boilerplate/go-api-helpers/service"

	"github.com/golang-devops/go-psexec/shared"
)

var (
	genpem  = flag.String("genpem", "", "The full path where to generate the pem file containing the private (and public) key")
	address = flag.String("address", ":62677", "The full host and port to listen on")
)

func main() {
	flag.Parse()

	if len(*genpem) > 0 {
		shared.GenerateKeyPairPemFile(*genpem)
		return
	}

	a := &app{}

	additionalArgs := []string{
		"-address",
		*address,
	}

	service.
		NewServiceRunnerBuilder("GoPsExec", a).
		WithOnStopHandler(a).
		WithAdditionalArguments(additionalArgs...).
		WithServiceUserName_AsCurrentUser().
		Run()
}
