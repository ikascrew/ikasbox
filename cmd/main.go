package main

import (
	"github.com/ikascrew/ikasbox"
)

func main() {

	err := ikasbox.Start()
	if err != nil {
		panic(err)
	}

}
