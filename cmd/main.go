package main

import (
	"github.com/ikascrew/ikabox"
)

func main() {

	err := ikabox.Start()
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
