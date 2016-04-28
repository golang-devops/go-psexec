# go-psexec
The plan is to have a replacement for psexec

## Getting Started

Very basic implementation to run a command remotely with no authentication yet.

### Get source
```
go get -u github.com/golang-devops/go-psexec
```

### Server

```
cd "%GOPATH%\src\github.com\golang-devops\go-psexec\server"
go build -o=server.exe
server.exe -allowed_public_keys_file "%UserProfile%/.config/go-psexec/server_allowed_public_keys" -server_pem "%UserProfile%/.config/go-psexec/server.pem"
```

#### Install server as service

```
cd "%GOPATH%\src\github.com\golang-devops\go-psexec\server"
go build -o=server.exe
server.exe -service install -allowed_public_keys_file "%UserProfile%/.config/go-psexec/server_allowed_public_keys" -server_pem "%UserProfile%/.config/go-psexec/server.pem"
```

If running under windows, the password needs to be set for the service.


### Client

#### Run the command-line interface (CLI)

```
cd "%GOPATH%\src\github.com\golang-devops\go-psexec\client\cmd"
go build -o=gopsexec-client.exe
gopsexec-client.exe -server "http://localhost:62677" -executor winshell ping 127.0.0.1 -n 6
```

#### Run from source

First check/replace variables in the `//Variables` section, **especially** the `clientPemFile`.

Create a `main.go` file with the content below and run `go run main.go`

```
package main

import (
    "log"
    "bufio"
    "fmt"

    "github.com/golang-devops/go-psexec/client"
    "github.com/golang-devops/go-psexec/shared"
    "github.com/golang-devops/go-psexec/shared/dtos"
)

func main() {
    //Variables
    clientPemFile := "/path/to/client.pem"
    serverUrl := "http://localhost:62677"
    executor := "winshell"
    exe := "ping"
    args := []string{"google.com", "-n", "6"}

    pvtKey, err := shared.ReadPemKey(clientPemFile)
    if err != nil {
        log.Fatalf("Cannot read client pem file, error: %s", err.Error())
    }

    c := client.New(pvtKey)

    session, err := c.RequestNewSession(serverUrl)
    if err != nil {
        log.Fatalf("Unable to create new session, error: %s", err.Error())
    }

    resp, err := session.StartExecRequest(&dtos.ExecDto{executor, exe, args})
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
}
```

## Known gotchas

There is some funny backslash escaping with the `interactive` mode. This is due to the way the library `github.com/nproc/parseargs-go` handles escaping of backslashes. So if you for instance run psexec with `psexec \\mymachine` it will escape the double backslash into a single.

## TODO

- Unit Tests
- Cross-platform support
- Running commands remotely and getting their feedback
- Authentication
- Perhaps support for LDAP
- Support for multiple executors (like using windows shell to "prepend" commands with `cmd /c`)
- Rework the server code to use `afero` file system instead of raw `os`. Much better for testability and abstraction