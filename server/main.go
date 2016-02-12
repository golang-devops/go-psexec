package main

import (
	"flag"
	"fmt"
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

	fmt.Println("............................................................................................................................")
	fmt.Println("............................................................................................................................")
	fmt.Println("............................................................................................................................")
	fmt.Println("TODO")
	fmt.Println("- Have a list of allowed client public-keys and check that in `checkPubKeyAllowed`")
	fmt.Println("- Encrypt the 'exec streamer feedback' with the session token.")
	fmt.Println("- More unit tests")
	fmt.Println("- Can perhaps make the token list (`tmpTokens`) a more persistent list using a DB or Redis.")
	fmt.Println("............................................................................................................................")
	fmt.Println("............................................................................................................................")
	fmt.Println("............................................................................................................................")

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
