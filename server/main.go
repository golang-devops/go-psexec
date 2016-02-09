package main

import (
	"flag"
	"github.com/zero-boilerplate/go-api-helpers/service"
)

var (
	address = flag.String("address", ":62677", "The full host and port to listen on")
)

func main() {
	flag.Parse()

	a := &app{}
	service.NewServiceRunnerBuilder("GoPsExec", a).Run()
}
