package main

import (
	"fmt"
	"os"

	"github.com/ikascrew/ikasbox"
)

func main() {

	err := ikasbox.Start()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}
