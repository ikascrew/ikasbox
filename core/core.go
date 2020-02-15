package core

import (
	"image"

	"gocv.io/x/gocv"
)

func ResizeImage(m gocv.Mat, w, h int) (gocv.Mat, error) {
	//dst := gocv.NewMatWithSize(w,h,gocv.MatTypeCV8U)
	dst := gocv.NewMat()
	gocv.Resize(m, &dst, image.Point{}, 0.5, 0.5, gocv.InterpolationDefault)
	return dst, nil
}

func WriteImage(f string, m gocv.Mat) error {
	gocv.IMWrite(f, m)
	return nil
}
