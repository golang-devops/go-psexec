package main

import (
	"flag"
	"github.com/zero-boilerplate/go-api-helpers/service"
)

const (
	Bearer     = "Bearer"
	SigningKey = "somethingsupersecret"
)

var (
	address = flag.String("address", ":62677", "The full host and port to listen on")
)

func main() {
	flag.Parse()

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
