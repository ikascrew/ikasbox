package main

import (
	"github.com/ikascrew/ikasbox"
)

func main() {

	err := ikasbox.Start()
	if err != nil {
		panic(err)
	}

	/*
		go func() {
			win, err := ikabox.NewWindow("test")
			if err != nil {
				panic(err)
			}
			defer win.Close()
			win.Load("test.mp4")

			for {
				win.Show()
				win.Wait()
			}
		}()
	*/
}
