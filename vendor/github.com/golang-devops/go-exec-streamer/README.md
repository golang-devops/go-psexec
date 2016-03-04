# go-exec-streamer
A golang exec streamer to make streaming exec.Command as simple as possible

## Usage

### Windows example

main.go:

```
package main

import (
    execstreamer "github.com/golang-devops/go-exec-streamer"
    "log"
    "os"
)

func main() {
    streamer, err := execstreamer.NewExecStreamerBuilder().
        ExecutorName("winshell").
        Exe("ping").
        Args("127.0.0.1", "-n", "6").
        StdoutWriter(os.Stdout).
        StderrWriter(os.Stderr).
        StdoutPrefix("OUT:").
        StderrPrefix("ERR:").
        AutoFlush().
        Build()

    if err != nil {
        log.Fatal(err)
    }

    err = streamer.ExecAndWait()
    if err != nil {
        log.Fatal(err)
    }
}
```

Now run it:

```
go run main.go
```

This should print out something like:

```
OUT:
OUT:Pinging 127.0.0.1 with 32 bytes of data:
OUT:Reply from 127.0.0.1: bytes=32 time<1ms TTL=128
OUT:Reply from 127.0.0.1: bytes=32 time<1ms TTL=128
OUT:Reply from 127.0.0.1: bytes=32 time<1ms TTL=128
OUT:Reply from 127.0.0.1: bytes=32 time<1ms TTL=128
OUT:Reply from 127.0.0.1: bytes=32 time<1ms TTL=128
OUT:Reply from 127.0.0.1: bytes=32 time<1ms TTL=128
OUT:
OUT:Ping statistics for 127.0.0.1:
OUT:    Packets: Sent = 6, Received = 6, Lost = 0 (0% loss),
OUT:Approximate round trip times in milli-seconds:
OUT:    Minimum = 0ms, Maximum = 0ms, Average = 0ms
```