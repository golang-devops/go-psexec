package main

import (
	"flag"
	"fmt"
	"github.com/zero-boilerplate/go-api-helpers/service"

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
		shared.GenerateKeyPairPemFile(*genpemFlag)
		return
	}

	if len(*genpubFromPemFlag) > 0 {
		shared.GeneratePublicKeyFromPemFlag(*genpubFromPemFlag)
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
		*addressFlag,
	}

	service.
		NewServiceRunnerBuilder("GoPsExec", a).
		WithOnStopHandler(a).
		WithAdditionalArguments(additionalArgs...).
		WithServiceUserName_AsCurrentUser().
		Run()
}
