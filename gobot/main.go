package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	oW := OmniWheelDriver{}
	oW.init(r, 5.6, 30, 30, 30)
	work := func() {
		//set spped
		oW.setMaxSpeed(15)
		fmt.Print("1")
		//Move forward one revolution
		//if err := ow.Move(2048); err != nil {
		//	fmt.Println(err)
		//}
		fmt.Print("1")
		oW.moveTo(20, 0) //10CM 30°
		//oW.rotate(90)
		oW.moveTo(20, 90) //10CM 30°
		//		oW.rotate(90)
		oW.moveTo(20, 180) //10CM 30°
		//		oW.rotate(90)
		oW.moveTo(20, -90) //10CM 30°
		//		oW.rotate(90)

		//oW.moveTo(10, 30) //10CM 30°
		//Move backward one revolution
		//	if err := stepper.Move(-2048); err != nil {
		//		fmt.Println(err)
		//	}
	}

	robot := gobot.NewRobot("stepperBot",
		[]gobot.Connection{r},
		[]gobot.Device{oW.stepper0, oW.stepper1, oW.stepper2, oW.stepper3},
		work,
	)

	robot.Start()
}
