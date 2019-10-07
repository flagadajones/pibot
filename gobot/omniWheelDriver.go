package main

import (
	"log"
	"math"
	"sync"
	"gobot.io/x/gobot/drivers/gpio"
)

const baseStepsPerRev = 2048
const maxSpeed = 9600
const maxAccel = 500
//const minPulse = 20

//var steppingTable [8]string = [8]string{"Full Step", "Half Step", "Quarter Step", "Eighth Step", "UNDEFINED Using Full Step", "UNDEFINED Using Full Step", "UNDEFINED Using Full Step", "Sisteenth Step"}
var stepsPerRev [8]int = [8]int{baseStepsPerRev, 2 * baseStepsPerRev, 4 * baseStepsPerRev, 8 * baseStepsPerRev, 0, 0, 0, 16 * baseStepsPerRev}

type OmniWheelDriver struct {
	stepper0 *gpio.StepperDriver
	stepper1 *gpio.StepperDriver
	stepper2 *gpio.StepperDriver
	stepper3 *gpio.StepperDriver
	index    int

	step   int
	sleep  bool
	enable bool
	maxSpeed      float64
	maxAccel      float64
	wheelDiameter float64
	baseDiameter  float64
}

func (omniW *OmniWheelDriver) init(driver gpio.DigitalWriter,wheelDiam float64, baseDiam float64, maxSpeed float64, maxAccel float64) {
	//1
	omniW.stepper0 = gpio.NewStepperDriver(driver, [4]string{"40", "38", "36", "32"}, gpio.StepperModes.DualPhaseStepping, 2048)
	//2
	omniW.stepper1 = gpio.NewStepperDriver(driver, [4]string{"37", "35", "33", "31"}, gpio.StepperModes.DualPhaseStepping, 2048)
	//3
	omniW.stepper2 = gpio.NewStepperDriver(driver, [4]string{"15", "13", "11", "7"}, gpio.StepperModes.DualPhaseStepping, 2048)
	//4
	omniW.stepper3 = gpio.NewStepperDriver(driver, [4]string{"22", "18", "16", "12"}, gpio.StepperModes.DualPhaseStepping, 2048)
	omniW.index = 0

	omniW.step = 0b001
	omniW.sleep = true
	omniW.enable = false

	omniW.maxSpeed = maxSpeed
	omniW.maxAccel = maxAccel
	omniW.wheelDiameter = wheelDiam
	omniW.baseDiameter = baseDiam

	/*
	   #  pinMode(ENABLE, OUTPUT)
	   #  pinMode(SLEEP, OUTPUT)
	   #  pinMode(RESET, OUTPUT)
	   #  pinMode(MS3, OUTPUT)
	   #  pinMode(MS2, OUTPUT)
	   #  pinMode(MS1, OUTPUT)
	*/
	//omniW.resetDriver()
	//omniW.setStep(omniW.step)
	
		omniW.setMaxSpeed(omniW.maxSpeed)
	//	omniW.setMaxAccel(omniW.maxAccel)
	//	omniW.setMinPulseWidth(minPulse)
	//	omniW.zeroLocation()
	//omniW.disableSteppers()
	//omniW.sleep()

}
/*
func (omniW *OmniWheelDriver) displayStatus() {

	log.Print("SLEEP: ")
	log.Print(omniW.sleep)
	log.Print("ENABLE: ")
	log.Print(omniW.enable)
	log.Print("OMNI WHEEL DIAMETER: ")
	log.Print(omniW.wheelDiameter)
	log.Print("cm")
	log.Print("BASE DIAMETER: ")
	log.Print(omniW.baseDiameter)
	log.Print("cm")
	log.Print("step = ")
	log.Print(steppingTable[omniW.step])
	log.Print("steps per rev = ")
	log.Print(stepsPerRev[omniW.step])
	log.Print("max speed = ")
	log.Print(omniW.maxSpeed)

}
*/
func (omniW *OmniWheelDriver) setBaseDiameter(diameter float64) {
	omniW.baseDiameter = diameter
	log.Print("setting base diameter to ")
	log.Print(omniW.baseDiameter)
	log.Print("cm")
}

func (omniW *OmniWheelDriver) setWheelDiameter(diameter float64) {

	omniW.wheelDiameter = diameter
	log.Print("setting omni wheel diameter to ")
	log.Print(omniW.wheelDiameter)
	log.Print("cm")
}

/*
   def set_step(omniW, step):
       omniW._step_ = step
       """   if (step & 0b100)
           digitalWrite(MS3, HIGH);
       else
           digitalWrite(MS3, LOW);

       if (step & 0b010)
           digitalWrite(MS2, HIGH);
       else
           digitalWrite(MS2, LOW);

       if (step & 0b001)
           digitalWrite(MS1, HIGH);
       else
           digitalWrite(MS1, LOW);
       """
       log.Print("step = ")
       log.Print(omniW.steppingTable[omniW._step_])
       log.Print("steps per rev = ")
       log.Print(omniW.stepsPerRev[omniW._step_])
*/

func (omniW *OmniWheelDriver) setMaxSpeed(speed float64) {

	omniW.maxSpeed = speed
	/*
		omniW.stepper0.setMaxSpeed(speed)
		omniW.stepper1.setMaxSpeed(speed)
		omniW.stepper2.setMaxSpeed(speed)
		omniW.stepper3.setMaxSpeed(speed)
	*/

	omniW.stepper0.SetSpeed(uint(speed))
	omniW.stepper1.SetSpeed(uint(speed))
	omniW.stepper2.SetSpeed(uint(speed))
	omniW.stepper3.SetSpeed(uint(speed))

	log.Print("max speed set to ")
	log.Print(speed)
}

/*
func (omniW *OmniWheelDriver) setMaxAccel(accel float64) {

	omniW.maxAccel = accel
	omniW.stepper0.setAcceleration(accel)
	omniW.stepper1.setAcceleration(accel)
	omniW.stepper2.setAcceleration(accel)
	omniW.stepper3.setAcceleration(accel)

	log.Print("max accel set to ")
	log.Print(accel)
}*/
/*
func (omniW *OmniWheelDriver) setMinPulseWidth(pulseWidth int) {

	omniW.stepper0.setMinPulseWidth(pulseWidth)
	omniW.stepper1.setMinPulseWidth(pulseWidth)
	omniW.stepper2.setMinPulseWidth(pulseWidth)
	omniW.stepper3.setMinPulseWidth(pulseWidth)
}
*/
/*

   def enable_steppers(omniW):
       #digitalWrite(ENABLE, LOW)
       omniW._enable_ = True


   def disable_steppers(omniW):
       #digitalWrite(ENABLE, HIGH);
       omniW._enable_ = False


   def sleep(omniW):
       #digitalWrite(SLEEP, LOW)
       omniW._sleep_ = True


   def wake(omniW):
       #digitalWrite(SLEEP, HIGH)
       omniW._sleep_ = False
*/
/*
func (omniW *OmniWheelDriver) zeroLocation() {

	omniW.stepper0.setCurrentPosition(0)
	omniW.stepper1.setCurrentPosition(0)
	omniW.stepper2.setCurrentPosition(0)
	omniW.stepper3.setCurrentPosition(0)
}
*/
/*
   def reset_driver(omniW):
       #digitalWrite(RESET, LOW)
       time.sleep(0.100)
       #digitalWrite(RESET, HIGH):
*/
/*
func (omniW *OmniWheelDriver) hardStopAll() {

	omniW.stepper0.setCurrentPosition(omniW.stepper0.targetPosition())
	omniW.stepper1.setCurrentPosition(omniW.stepper1.targetPosition())
	omniW.stepper2.setCurrentPosition(omniW.stepper2.targetPosition())
	omniW.stepper3.setCurrentPosition(omniW.stepper3.targetPosition())
}
func (omniW *OmniWheelDriver) run() {
	omniW.stepper0.run()
	omniW.stepper1.run()
	omniW.stepper2.run()
	omniW.stepper3.run()
}*/
func (omniW *OmniWheelDriver) stopAll() {
	omniW.stepper0.Halt()
	omniW.stepper1.Halt()
	omniW.stepper2.Halt()
	omniW.stepper3.Halt()
}

func (omniW *OmniWheelDriver) rotate(deg float64) {
	var wg sync.WaitGroup
	//omniW.setMaxAccel(omniW.maxAccel)
	baseCircumference := math.Pi * omniW.baseDiameter
	wheelCircumference := math.Pi * omniW.wheelDiameter
	cmPerStep := wheelCircumference / float64(baseStepsPerRev/*stepsPerRev[omniW.step]*/)
	steps := ((deg / 360.0) * baseCircumference) / cmPerStep
	wg.Add(4)
 go	func() {omniW.stepper0.Move(int(steps))
	defer wg.Done()}()
go	func() {omniW.stepper1.Move(int(steps))
defer wg.Done()}()
go	func() {omniW.stepper2.Move(int(steps))
defer wg.Done()}()
go	func() {omniW.stepper3.Move(int(steps))
defer wg.Done()}()
wg.Wait()
}

func (omniW *OmniWheelDriver) moveTo(r float64, theta float64) {
	var wg sync.WaitGroup
	thetaRad := math.Pi / 180.0 * theta
	x := r * math.Cos(thetaRad)
	y := r * math.Sin(thetaRad)
	wheelCircumference := math.Pi * omniW.wheelDiameter
	cmPerStep := wheelCircumference / float64(baseStepsPerRev/*stepsPerRev[omniW.step]*/)
	xSteps := x / cmPerStep
	ySteps := y / cmPerStep

	if math.Abs(xSteps) > math.Abs(ySteps) {
	//	slowPercent := math.Abs(ySteps / xSteps)
	//	slowSpeed := slowPercent * omniW.maxAccel
	//	omniW.stepper0.setAcceleration(omniW.maxAccel)
	//	omniW.stepper1.setAcceleration(omniW.maxAccel)
	//	omniW.stepper2.setAcceleration(slowSpeed)
	//	omniW.stepper3.setAcceleration(slowSpeed)
	} else if math.Abs(xSteps) < math.Abs(ySteps) {
	//	slowPercent := math.Abs(xSteps / ySteps)
	//	slowSpeed := slowPercent * omniW.maxAccel
	//	omniW.stepper0.setAcceleration(slowSpeed)
	//	omniW.stepper1.setAcceleration(slowSpeed)
	//	omniW.stepper2.setAcceleration(omniW.maxAccel)
	//	omniW.stepper3.setAcceleration(omniW.maxAccel)
	} else {
	//	omniW.setMaxAccel(omniW.maxAccel)
	}
	wg.Add(4)
	go func() { omniW.stepper0.Move(int(xSteps))
	defer wg.Done()}()
	go func() { omniW.stepper1.Move(int(-xSteps))
	defer wg.Done()}()
	go func() { omniW.stepper2.Move(int(ySteps))
	defer wg.Done()}()
	go func() {omniW.stepper3.Move(int(-ySteps))
	defer wg.Done()}()

	wg.Wait()
}
