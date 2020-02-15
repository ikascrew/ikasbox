package core

import (
	"fmt"

	"gocv.io/x/gocv"
)

type Window struct {
	owner   *gocv.Window
	capture *gocv.VideoCapture
	mat     gocv.Mat
	frames  float64
}

func NewWindow(name string) (*Window, error) {
	win := gocv.NewWindow(name)
	w := &Window{
		owner:   win,
		capture: nil,
		mat:     gocv.NewMat(),
		frames:  0.0,
	}
	return w, nil
}

func (w *Window) Load(file string) error {
	w.CloseCapture()
	cap, err := gocv.VideoCaptureFile(file)
	if err != nil {
		return err
	}
	w.capture = cap
	w.frames = cap.Get(gocv.VideoCaptureFrameCount)
	return nil
}

func (w *Window) IsOpened() bool {
	if w.capture != nil {
		if w.capture.IsOpened() {
			return true
		}
	}
	return false
}

func (w *Window) Show() error {

	if !w.IsOpened() {
		return fmt.Errorf("Capture Not Open")
	}

	flag := w.capture.Read(&w.mat)
	if flag {
		//return fmt.Errorf("Capture Error?")
	}
	w.owner.IMShow(w.mat)
	pos := w.capture.Get(gocv.VideoCapturePosFrames)
	if w.frames == pos {
		w.capture.Set(gocv.VideoCapturePosFrames, 1)
	}
	return nil
}

func (w *Window) Wait() {
	gocv.WaitKey(33)
}

func (w *Window) CloseCapture() bool {
	if w.IsOpened() {
		w.capture.Close()
		return true
	}
	return false
}

func (w *Window) Close() {
	w.CloseCapture()
	w.owner.Close()
}
