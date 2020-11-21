package main

import (
	"gocv.io/x/gocv"
	"image"
	"log"
	"time"
)

func CameraDemo() {
	webcam, _ := gocv.OpenVideoCapture(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}
}
func Crop(imgPath string, savePath string, rect image.Rectangle) {
	t1:=time.Now()  //获取本地现在时间
	img := gocv.IMRead(imgPath, gocv.IMReadColor)
	defer img.Close()
	img = img.Region(rect)
	gocv.IMWrite(savePath, img)
	log.Printf("Crop spend: %s", time.Now().Sub(t1))
}
