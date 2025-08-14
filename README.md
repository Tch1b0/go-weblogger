# go-weblogger

Expose your terminal output to the web

## usage

```go
package main

import (
	"fmt"
	"time"

    weblogger "github.com/Tch1b0/go-weblogger"
)

func main() {
    wi := weblogger.NewWebInterface("localhost", 3000)
    go wi.Serve()

    for {
        fmt.Fprintln(wi.Writer, "This is a message to the webterminal")
        fmt.Fprintln(wi.Writer, "This is a second message")
        fmt.Fprintln(wi.Writer, "You can see, the messages are \n displayed sequentially top-down")

        time.Sleep(3 * time.Second)
    }
}
```

![example output](./example.png)
