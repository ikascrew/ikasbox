package main

import (
	"fmt"
	"os"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("%+v")
		os.Exit(1)
	}
	fmt.Println("Success.")
}

func run() error {

	f := "sample.jpg"

	return nil
}
