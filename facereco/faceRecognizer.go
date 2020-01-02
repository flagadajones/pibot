package facereco

import (
	"image"
	"path"
	"runtime"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"gocv.io/x/gocv"
)

// FaceRecognizer ...
type FaceRecognizer struct {
	cascade string
	camera  *opencv.CameraDriver
	img     atomic.Value
	Devices []gobot.Device
	window  *gocv.Window
	faces   *[]image.Rectangle
}

// GetFaces ...
func (faceR *FaceRecognizer) GetFaces() *[]image.Rectangle {
	return faceR.faces

}

// Init ...
func (faceR *FaceRecognizer) Init() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	//window = opencv.NewWindowDriver()
	faceR.camera = opencv.NewCameraDriver(0)

	mat := gocv.NewMat()
	faceR.img.Store(mat)
	faceR.camera.On(opencv.Frame, func(data interface{}) {
		j := data.(gocv.Mat)
		faceR.img.Store(j)
	})
	gobot.Every(30*time.Millisecond, func() {
		j := faceR.img.Load().(gocv.Mat)
		if j.Empty() {
			return
		}
		//	log.Println("size")
		//	log.Println(j.Size())
		faces := opencv.DetectObjects(cascade, j)
		faceR.faces = &faces
		//log.Println(window)
		opencv.DrawRectangles(j, faces, 0, 255, 0, 5)
		//i = j
		//	gocv.IMWrite("/Users/grua341/go/src/github.com/flagadajones/pibot/save.jpeg", j)
		faceR.img.Store(j)

		//	faceR.Run()

	})
	//faceR.window = gocv.NewWindow("Hello")
	faceR.Devices = []gobot.Device{ /*window,*/ faceR.camera}
}

// Run ...
func (faceR *FaceRecognizer) Run() {

	for {
		i := faceR.img.Load().(gocv.Mat)

		if len(i.Size()) != 0 {

			faceR.window.IMShow(i)
			faceR.window.WaitKey(1)
		}
		//time.Sleep(1 * time.Second)
	}
}
