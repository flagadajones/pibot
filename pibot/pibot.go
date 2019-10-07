package main

import (
	"log"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/flagadajones/pibot/facereco"
	"github.com/flagadajones/pibot/visage"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

var devices = []gobot.Device{}

var robot *gobot.Robot
var facerec *facereco.FaceRecognizer

func run() {
	log.Print("eee")
	robot.Start(true)

	//	work := func() {
	//		v.Run()
	//	oW.moveTo(20, 0) //10CM 30°
	//	oW.moveTo(20, 90) //10CM 30°
	//	oW.moveTo(20, 180) //10CM 30°
	//	oW.moveTo(20, -90) //10CM 30°
	//	}
	//facereco.Run()
}
func init() {

}

func main() {
	r := raspi.NewAdaptor()
	facerec = &facereco.FaceRecognizer{}
	facerec.Init()

	v := visage.Visage{}
	v.Init(&visage.CapSize{W: 1280, H: 720})
	work := func() {
		//facerec.Run()
		v.Run()

		//	oW.moveTo(20, 0) //10CM 30°
		//	oW.moveTo(20, 90) //10CM 30°
		//	oW.moveTo(20, 180) //10CM 30°
		//	oW.moveTo(20, -90) //10CM 30°
	}
	devices = append(devices, facerec.Devices...)
	robot = gobot.NewRobot("stepperBot",
		[]gobot.Connection{r},
		devices,
		work,
	)
	gobot.Every(3000*time.Millisecond, func() {
		v.Blink()
	})
	gobot.Every(10*time.Millisecond, func() {
		faces := facerec.GetFaces()
		if faces != nil && len(*faces) > 0 {
			log.Println("face")
			log.Println((*faces)[0])
			log.Println("f1", 1280-(*faces)[0].Min.X, 720-(*faces)[0].Min.Y)
			w := (*faces)[0].Max.X - (*faces)[0].Min.X
			h := (*faces)[0].Max.Y - (*faces)[0].Min.Y
			v.Draw(&visage.Cible{X: 1280 - (*faces)[0].Max.X, Y: 720 - (*faces)[0].Min.Y, W: w, H: h})
		}
	})

	pixelgl.Run(run)

}
