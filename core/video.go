package core

import (
	"bytes"
	"fmt"

	"gocv.io/x/gocv"
)

type Video struct {
	owner  *gocv.VideoCapture
	name   string
	Width  int
	Height int
	FPS    float64
	FOURCC float64
	Frames int
}

func NewVideo(f string) (*Video, error) {

	cap, err := gocv.VideoCaptureFile(f)
	if err != nil {
		return nil, err
	}

	width := cap.Get(gocv.VideoCaptureFrameWidth)
	height := cap.Get(gocv.VideoCaptureFrameHeight)
	fps := cap.Get(gocv.VideoCaptureFPS)
	frames := cap.Get(gocv.VideoCaptureFrameCount)
	fourcc := cap.Get(gocv.VideoCaptureFOURCC)

	v := &Video{
		owner:  cap,
		name:   f,
		Width:  int(width),
		Height: int(height),
		FPS:    fps,
		FOURCC: fourcc,
		Frames: int(frames),
	}
	return v, nil
}

func (v *Video) Close() {
	err := v.owner.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (v *Video) GetImage(f float64) (*gocv.Mat, error) {
	if v == nil || !v.owner.IsOpened() {
		return nil, fmt.Errorf("Capture Open Error")
	}
	m := gocv.NewMat()
	v.owner.Set(gocv.VideoCapturePosFrames, f)
	v.owner.Read(&m)
	return &m, nil
}

func (v *Video) String() string {
	w := bytes.NewBufferString("")
	fmt.Fprintf(w, "File   :[%s] {\n", v.name)
	fmt.Fprintf(w, "  Width  :[%d]\n", v.Width)
	fmt.Fprintf(w, "  Height :[%d]\n", v.Height)
	fmt.Fprintf(w, "  FPS    :[%f]\n", v.FPS)
	fmt.Fprintf(w, "  Frames :[%d]\n", v.Frames)
	fmt.Fprintf(w, "  FOURCC :[%f]\n", v.FOURCC)
	fmt.Fprintf(w, "}\n")
	return w.String()
}
