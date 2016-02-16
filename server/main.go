package main

import (
	"flag"
	"github.com/zero-boilerplate/go-api-helpers/service"
	"log"

	"github.com/golang-devops/go-psexec/shared"
)

var (
	//These flags is not run as service but will exit after completion
	genpemFlag        = flag.String("genpem", "", "The full path where to generate the pem file containing the private (and public) key")
	genpubFromPemFlag = flag.String("pub_from_pem", "", "Generate the public key from the input pem file")
)

var (
	addressFlag               = flag.String("address", ":62677", "The full host and port to listen on")
	allowedPublicKeysFileFlag = flag.String("allowed_public_keys_file", "", "The path to the allowed public keys file")
)

func main() {
	flag.Parse()

	if len(*genpemFlag) > 0 {
		err := shared.GenerateKeyPairPemFile(*genpemFlag)
		if err != nil {
			log.Fatalf("Unable to generate key pair pem file, error: %s", err.Error())
		}
		return
	}

	if len(*genpubFromPemFlag) > 0 {
		err := shared.PrintPemFilePublicKeyAsHex(*genpubFromPemFlag)
		if err != nil {
			log.Fatalf("Unable to generate public key from pem file, error: %s", err.Error())
		}
		return
	}

	if len(*allowedPublicKeysFileFlag) == 0 {
		log.Fatalln("No allowed public keys file specified, no keys will be allowed. Exiting server.")
	}

	a := &app{}

	additionalArgs := []string{
		"-address",
		*addressFlag,
		"-allowed_public_keys_file",
		*allowedPublicKeysFileFlag,
	}

	service.
		NewServiceRunnerBuilder("GoPsExec", a).
		WithOnStopHandler(a).
		WithAdditionalArguments(additionalArgs...).
		WithServiceUserName_AsCurrentUser().
		Run()
}
