package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Printf("os.Stdout = %p, type = %T\n", os.Stdout, os.Stdout)

	wi := NewWebInterface("localhost", 3000)
	go wi.Start()

	fmt.Printf("os.Stdout = %p, type = %T\n", os.Stdout, os.Stdout)

	for {

		fmt.Fprint(wi.Writer, "This is a message to the webterminal")
		fmt.Fprint(wi.Writer, "and another message")
		fmt.Fprint(wi.Writer, "You can see, the messages are displayed sequentially top-down")

		fmt.Printf("os.Stdout = %p, type = %T\n", os.Stdout, os.Stdout)
		time.Sleep(5 * time.Second)
	}
}
