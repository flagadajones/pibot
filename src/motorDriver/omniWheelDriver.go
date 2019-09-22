package main

import (
	"log"
	"math"
)

const baseStepsPerRev = 200
const maxSpeed = 9600
const maxAccel = 500
const minPulse = 20

var steppingTable [8]string = [8]string{"Full Step", "Half Step", "Quarter Step", "Eighth Step", "UNDEFINED Using Full Step", "UNDEFINED Using Full Step", "UNDEFINED Using Full Step", "Sisteenth Step"}
var stepsPerRev [8]int = [8]int{baseStepsPerRev, 2 * baseStepsPerRev, 4 * baseStepsPerRev, 8 * baseStepsPerRev, 0, 0, 0, 16 * baseStepsPerRev}

type OmniWheelDriver struct {
	stepper0 AccelStepper
	stepper1 AccelStepper
	stepper2 AccelStepper
	stepper3 AccelStepper
	index    int

	step   int
	sleep  bool
	enable bool

	maxSpeed      float64
	maxAccel      float64
	wheelDiameter float64
	baseDiameter  float64
}

func (omniW *OmniWheelDriver) init(wheelDiam float64, baseDiam float64, maxSpeed float64, maxAccel float64) {
	//omniW.stepper0 = AccelStepper{}
	omniW.stepper0.init(FULL4WIRE, 37,35, 33, 31, true)
	//omniW.stepper1 = AccelStepper{}
	omniW.stepper1.init(FULL4WIRE, 1, 2, 3, 4, false)
	//omniW.stepper2 = AccelStepper{}
	omniW.stepper2.init(FULL4WIRE, 1, 2, 3, 4, false)
	//omniW.stepper3 = AccelStepper{}
	omniW.stepper3.init(FULL4WIRE, 1, 2, 3, 4, false)
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
	omniW.setMaxAccel(omniW.maxAccel)
	omniW.setMinPulseWidth(minPulse)
	omniW.zeroLocation()
	//omniW.disableSteppers()
	//omniW.sleep()

}

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
	omniW.stepper0.setMaxSpeed(speed)
	omniW.stepper1.setMaxSpeed(speed)
	omniW.stepper2.setMaxSpeed(speed)
	omniW.stepper3.setMaxSpeed(speed)

	log.Print("max speed set to ")
	log.Print(speed)
}

func (omniW *OmniWheelDriver) setMaxAccel(accel float64) {

	omniW.maxAccel = accel
	omniW.stepper0.setAcceleration(accel)
	omniW.stepper1.setAcceleration(accel)
	omniW.stepper2.setAcceleration(accel)
	omniW.stepper3.setAcceleration(accel)

	log.Print("max accel set to ")
	log.Print(accel)
}
func (omniW *OmniWheelDriver) setMinPulseWidth(pulseWidth int) {

	omniW.stepper0.setMinPulseWidth(pulseWidth)
	omniW.stepper1.setMinPulseWidth(pulseWidth)
	omniW.stepper2.setMinPulseWidth(pulseWidth)
	omniW.stepper3.setMinPulseWidth(pulseWidth)
}

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
func (omniW *OmniWheelDriver) zeroLocation() {

	omniW.stepper0.setCurrentPosition(0)
	omniW.stepper1.setCurrentPosition(0)
	omniW.stepper2.setCurrentPosition(0)
	omniW.stepper3.setCurrentPosition(0)
}

/*
   def reset_driver(omniW):
       #digitalWrite(RESET, LOW)
       time.sleep(0.100)
       #digitalWrite(RESET, HIGH):
*/
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
}
func (omniW *OmniWheelDriver) stopAll() {
	omniW.stepper0.stop()
	omniW.stepper1.stop()
	omniW.stepper2.stop()
	omniW.stepper3.stop()
}

func (omniW *OmniWheelDriver) rotate(deg float64) {

	omniW.setMaxAccel(omniW.maxAccel)
	baseCircumference := math.Pi * omniW.baseDiameter
	wheelCircumference := math.Pi * omniW.wheelDiameter
	cmPerStep := wheelCircumference / float64(stepsPerRev[omniW.step])
	steps := ((deg / 360.0) * baseCircumference) / cmPerStep

	omniW.stepper0.move(int(steps))
	omniW.stepper1.move(int(steps))
	omniW.stepper2.move(int(steps))
	omniW.stepper3.move(int(steps))
}

func (omniW *OmniWheelDriver) moveTo(r float64, theta float64) {

	thetaRad := math.Pi / 180.0 * theta
	x := r * math.Cos(thetaRad)
	y := r * math.Sin(thetaRad)
	wheelCircumference := math.Pi * omniW.wheelDiameter
	cmPerStep := wheelCircumference / float64(stepsPerRev[omniW.step])
	xSteps := x / cmPerStep
	ySteps := y / cmPerStep

	if math.Abs(xSteps) > math.Abs(ySteps) {
		slowPercent := math.Abs(ySteps / xSteps)
		slowSpeed := slowPercent * omniW.maxAccel
		omniW.stepper0.setAcceleration(omniW.maxAccel)
		omniW.stepper1.setAcceleration(omniW.maxAccel)
		omniW.stepper2.setAcceleration(slowSpeed)
		omniW.stepper3.setAcceleration(slowSpeed)
	} else if math.Abs(xSteps) < math.Abs(ySteps) {
		slowPercent := math.Abs(xSteps / ySteps)
		slowSpeed := slowPercent * omniW.maxAccel
		omniW.stepper0.setAcceleration(slowSpeed)
		omniW.stepper1.setAcceleration(slowSpeed)
		omniW.stepper2.setAcceleration(omniW.maxAccel)
		omniW.stepper3.setAcceleration(omniW.maxAccel)
	} else {
		omniW.setMaxAccel(omniW.maxAccel)
	}
	omniW.stepper0.move(int(xSteps))
	omniW.stepper1.move(int(-xSteps))
	omniW.stepper2.move(int(ySteps))
	omniW.stepper3.move(int(-ySteps))
}
