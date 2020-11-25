package main

import (
	logger "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
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

func Crop(img gocv.Mat, rect image.Rectangle, savePath string) {
	logger.Debug(savePath)
	t1:=time.Now()  //获取本地现在时间
	img = img.Region(rect)
	gocv.IMWrite(savePath, img)
	logger.Debug(time.Now().Sub(t1))
}

func CropByPath(imgPath string, savePath string, rect image.Rectangle) {
	t1:=time.Now()  //获取本地现在时间
	img := gocv.IMRead(imgPath, gocv.IMReadColor)
	defer img.Close()
	img = img.Region(rect)
	gocv.IMWrite(savePath, img)
	logger.Debug(time.Now().Sub(t1))
}
