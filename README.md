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
server.exe
```

#### Install server as service

```
cd "%GOPATH%\src\github.com\golang-devops\go-psexec\server"
go build -o=server.exe
server.exe -service install
```

If running under windows, the password needs to be set for the service.


Client:

```
cd "%GOPATH%\src\github.com\golang-devops\go-psexec\client"
go build -o=client.exe
client.exe -server "http://localhost:62677" -executor winshell ping 127.0.0.1 -n 6
```

## TODO

- Unit Tests
- Cross-platform support
- Running commands remotely and getting their feedback
- Authentication
- Perhaps support for LDAP
- Support for multiple executors (like using windows shell to "prepend" commands with `cmd /c`)