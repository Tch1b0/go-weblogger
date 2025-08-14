package main

import (
	"fmt"
	"time"
)

func main() {
	wi := NewWebInterface("localhost", 3000)
	go wi.Serve()

	time.Sleep(3 * time.Second)

	fmt.Fprint(wi.Writer, "This is a message to the webterminal")
	fmt.Fprint(wi.Writer, "and another message")
	fmt.Fprint(wi.Writer, "You can see, the messages are \n displayed sequentially top-down")

	time.Sleep(3 * time.Second)
}
