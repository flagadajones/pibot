package main

import (
	"math"
	"time"

	//rpi "github.com/nathan-osman/go-rpigpio"
	rpi "github.com/nahsh/go-rpigpio"
)

func constrain(val, minVal, maxVal float64) float64 {

	return math.Min(maxVal, math.Max(minVal, val))

}

//
const (
	DirectionCW  = 0 // Counter-Clockwise
	DirectionCCW = 1 // Clockwise

)

//
const (
	//class MotorInterfaceType(Enum):
	FUNCTION  = 0 // //< Use the functional interface, implementing your own driver functions (internal use only)
	DRIVER    = 1 // //< Stepper Driver, 2 driver pins required
	FULL2WIRE = 2 // //< 2 wire stepper, 2 motor pins required
	FULL3WIRE = 3 // //< 3 wire stepper, such as HDD spindle, 3 motor pins required
	FULL4WIRE = 4 ////< 4 wire full stepper, 4 motor pins required
	HALF3WIRE = 6 ////< 3 wire half stepper, such as HDD spindle, 3 motor pins required
	HALF4WIRE = 8 ////< 4 wire half stepper, 4 motor pins required
)

// AccelStepper  accelStepper
type AccelStepper struct {
	pin         [4]int
	pinInverted [4]int
	pinGPIO     [4]*rpi.Pin

	interf        int
	currentPos    int
	targetPos     int
	speed         float64
	maxSpeed      float64
	acceleration  float64
	sqrtTwoa      float64
	stepInterval  int
	minPulseWidth int
	enablePin     byte
	lastStepTime  int

	enableInverted bool

	n         float64 //int
	c0        float64
	cn        float64
	cmin      float64
	direction int
}

func (stepper *AccelStepper) forward()  {}
func (stepper *AccelStepper) backward() {}
func (stepper *AccelStepper) init(interf, pin1, pin2, pin3, pin4 int, enable bool) {

	stepper.interf = interf
	stepper.currentPos = 0
	stepper.targetPos = 0
	stepper.speed = 0.0
	stepper.maxSpeed = 1.0
	stepper.acceleration = 0.0
	stepper.sqrtTwoa = 1.0
	stepper.stepInterval = 0
	stepper.minPulseWidth = 1
	stepper.enablePin = 0xff
	stepper.lastStepTime = 0
	stepper.pin[0] = pin1
	stepper.pin[1] = pin2
	stepper.pin[2] = pin3
	stepper.pin[3] = pin4
	stepper.enableInverted = false

	//stepper.forward=lambda *_, **_: None
	//stepper.backward=lambda *_, **_: None

	//// NEW
	stepper.n = 0
	stepper.c0 = 0.0
	stepper.cn = 0.0
	stepper.cmin = 1.0
	stepper.direction = DirectionCCW

	for i := 0; i < 4; i++ {
		stepper.pinInverted[i] = 0
	}

	if enable {
		stepper.enableOutputs()
	}
	//// Some reasonable default
	stepper.setAcceleration(1)

}

func (stepper *AccelStepper) setSpeed(speed float64) {

	if speed == stepper.speed {
		return
	}
	speed = constrain(speed, -stepper.maxSpeed, stepper.maxSpeed)
	if speed == 0.0 {
		stepper.stepInterval = 0
	} else {
		stepper.stepInterval = int(math.Abs(1000000.0 / speed))

		if speed > 0.0 {
			stepper.direction = DirectionCW
		} else {
			stepper.direction = DirectionCCW
		}
	}
	stepper.speed = speed
}

func (stepper *AccelStepper) getSpeed() float64 {
	return stepper.speed
}

// Subclasses can override
func (stepper *AccelStepper) step(step int) {
	switch stepper.interf {
	case FUNCTION:
		stepper.step0(step)
	case DRIVER:
		stepper.step1(step)
	case FULL2WIRE:
		stepper.step2(step)
	case FULL3WIRE:
		stepper.step3(step)
	case FULL4WIRE:
		stepper.step4(step)
	case HALF3WIRE:
		stepper.step6(step)
	case HALF4WIRE:
		stepper.step8(step)
	}

}

// You might want to override this to implement eg serial output
// bit 0 of the mask corresponds to _pin[0]
// bit 1 of the mask corresponds to _pin[1]
// ....
func (stepper *AccelStepper) setOutputPins(mask byte) {
	numpins := 2
	if stepper.interf == FULL4WIRE || stepper.interf == HALF4WIRE {
		numpins = 4
	} else if stepper.interf == FULL3WIRE || stepper.interf == HALF3WIRE {
		numpins = 3
	}

	for i := 0; i < numpins; i++ {
		if mask&(1<<i) != 0 {
			stepper.pinGPIO[i].Write(rpi.HIGH /*^ stepper.pinInverted[i]*/)
			//GPIO.output(stepper.pin[i], GPIO.HIGH^stepper.pinInverted[i])
		} else {
			stepper.pinGPIO[i].Write(rpi.LOW /*^ stepper.pinInverted[i]*/)
			//GPIO.output(stepper.pin[i], GPIO.LOW^stepper.pinInverted[i]) // 	digitalWrite(_pin[i], (mask & (1 << i)) ? (HIGH ^ _pinInverted[i]) : (LOW ^ _pinInverted[i]));
		}
	}
}

// 0 pin step function (ie for functional usage)
func (stepper *AccelStepper) step0(step int) {
	//(void)(step); // Unused
	if stepper.speed > 0 {
		stepper.forward()
	} else {
		stepper.backward()
	}
}

// 1 pin step function (ie for stepper drivers)
// This is passed the current step number (0 to 7)
// Subclasses can override
func (stepper *AccelStepper) step1(step int) {
	// (void)(step); // Unused
	// _pin[0] is step, _pin[1] is direction
	if stepper.direction == DirectionCW {
		stepper.setOutputPins(0b10) // Set direction first else get rogue pulses
		stepper.setOutputPins(0b11) // step HIGH
	} else {
		stepper.setOutputPins(0b00) // Set direction first else get rogue pulses
		stepper.setOutputPins(0b01) // step HIGH
	}
	// Caution 200ns setup time
	// Delay the minimum allowed pulse width
	time.Sleep(time.Duration(stepper.minPulseWidth) * time.Microsecond) //// delayMicroseconds(stepper.minPulseWidth)
	if stepper.direction == DirectionCW {
		stepper.setOutputPins(0b10) // step LOW
	} else {
		stepper.setOutputPins(0b00) // step LOW
	}
}

// 2 pin step function
// This is passed the current step number (0 to 7)
// Subclasses can override
func (stepper *AccelStepper) step2(step int) {

	switch step & 0x3 {
	case 0: // 01
		stepper.setOutputPins(0b10)
	case 1: // 11
		stepper.setOutputPins(0b11)
	case 2: // 10
		stepper.setOutputPins(0b01)
	case 3: // 00
		stepper.setOutputPins(0b00)
	}
}

// 3 pin step function
// This is passed the current step number (0 to 7)
// Subclasses can override
func (stepper *AccelStepper) step3(step int) {

	switch step % 3 {
	case 0: // 100
		stepper.setOutputPins(0b100)
	case 1: // 001
		stepper.setOutputPins(0b001)
	case 2: // 010
		stepper.setOutputPins(0b010)

	}
}

// 4 pin step function for half stepper
// This is passed the current step number (0 to 7)
// Subclasses can override
func (stepper *AccelStepper) step4(step int) {

	switch step & 0x3 {
	case 0: // 1010
		stepper.setOutputPins(0b0101)
	case 1: // 0110
		stepper.setOutputPins(0b0110)
	case 2: // 0101
		stepper.setOutputPins(0b1010)
	case 3: // 1001
		stepper.setOutputPins(0b1001)
	}
}

// 3 pin half step function
// This is passed the current step number (0 to 7)
// Subclasses can override
func (stepper *AccelStepper) step6(step int) {

	switch step % 6 {
	case 0: // 100
		stepper.setOutputPins(0b100)
	case 1: // 101
		stepper.setOutputPins(0b101)
	case 2: // 001
		stepper.setOutputPins(0b001)
	case 3: // 011
		stepper.setOutputPins(0b011)
	case 4: // 010
		stepper.setOutputPins(0b010)
	case 5: // 110
		stepper.setOutputPins(0b110)
	}
}

// 4 pin half step function
// This is passed the current step number (0 to 7)
// Subclasses can override
func (stepper *AccelStepper) step8(step int) {

	switch step & 0x7 {
	case 0: // 1000
		stepper.setOutputPins(0b0001)
	case 1: // 1010
		stepper.setOutputPins(0b0101)
	case 2: // 0010
		stepper.setOutputPins(0b0100)
	case 3: // 0110
		stepper.setOutputPins(0b0110)
	case 4: // 0100
		stepper.setOutputPins(0b0010)
	case 5: // 0101
		stepper.setOutputPins(0b1010)
	case 6: // 0001
		stepper.setOutputPins(0b1000)
	case 7: // 1001
		stepper.setOutputPins(0b1001)
	}
}

// Implements steps according to the current step interval
// You must call this at least once per step
// returns true if a step occurred
func (stepper *AccelStepper) runSpeed() bool {
	// Dont do anything unless we actually have a step interval
	if stepper.stepInterval == 0 {
		return false
	}

	timer := int(time.Now().UnixNano() / 1000)
	//FIXME
	//	timer := datetime.now().microseconds()
	if timer-stepper.lastStepTime >= stepper.stepInterval {
		if stepper.direction == DirectionCW {
			// Clockwise
			stepper.currentPos++
		} else {
			// Anticlockwise
			stepper.currentPos--
		}
		stepper.step(stepper.currentPos)

		stepper.lastStepTime = timer // Caution: does not account for costs in step()

		return true
	} else {
		return false
	}
}

// Run the motor to implement speed and acceleration in order to proceed to the target position
// You must call this at least once per step, preferably in your main loop
// If the motor is in the desired position, the cost is very small
// returns true if the motor is still running to the target position.
func (stepper *AccelStepper) run() bool {
	if stepper.runSpeed() {
		stepper.computeNewSpeed()
	}
	return stepper.speed != 0.0 || stepper.distanceToGo() != 0
}

func (stepper *AccelStepper) setMaxSpeed(speed float64) {
	if speed < 0.0 {
		speed = -speed
	}
	if stepper.maxSpeed != speed {
		stepper.maxSpeed = speed
		stepper.cmin = 1000000.0 / speed
		// Recompute _n from current speed and adjust speed if accelerating or cruising
		if stepper.n > 0 {
			stepper.n = ((stepper.speed * stepper.speed) / (2.0 * stepper.acceleration)) // Equation 16
			stepper.computeNewSpeed()
		}
	}
}
func (stepper *AccelStepper) getMaxSpeed() float64 {
	return stepper.maxSpeed
}

func (stepper *AccelStepper) enableOutputs() {

	if stepper.interf == 0 {
		return
	}

	p0, err := rpi.OpenPin(stepper.pin[0], rpi.OUT)
	stepper.pinGPIO[0] = p0
	if err != nil {
		panic(err)
	}
	defer stepper.pinGPIO[0].Close()

	p1, err := rpi.OpenPin(stepper.pin[1], rpi.OUT)
	stepper.pinGPIO[1] = p1
	if err != nil {
		panic(err)
	}
	defer stepper.pinGPIO[1].Close()

	if stepper.interf == FULL4WIRE || stepper.interf == HALF4WIRE {
		p2, err := rpi.OpenPin(stepper.pin[2], rpi.OUT)
		stepper.pinGPIO[2] = p2
		if err != nil {
			panic(err)
		}
		defer stepper.pinGPIO[2].Close()

		p3, err := rpi.OpenPin(stepper.pin[3], rpi.OUT)
		stepper.pinGPIO[3] = p3
		if err != nil {
			panic(err)
		}
		defer stepper.pinGPIO[3].Close()

	} else if stepper.interf == FULL3WIRE || stepper.interf == HALF3WIRE {
		p2, err := rpi.OpenPin(stepper.pin[2], rpi.OUT)
		stepper.pinGPIO[2] = p2
		if err != nil {
			panic(err)
		}
		defer stepper.pinGPIO[2].Close()

	}
	//	if stepper.enablePin != 0xff {
	//		GPIO.setup(stepper.enablePin, GPIO.OUT)
	//		GPIO.output(stepper.enablePin, GPIO.HIGH^stepper.enableInverted) //   digitalWrite(_enablePin, GPIO.HIGH ^ stepper.enableInverted);
	//	}
}

func (stepper *AccelStepper) setAcceleration(acceleration float64) {

	if acceleration == 0.0 {
		return
	}
	if acceleration < 0.0 {
		acceleration = -acceleration
	}
	if stepper.acceleration != acceleration {
		// Recompute _n per Equation 17
		stepper.n = stepper.n * (stepper.acceleration / acceleration)
		// New c0 per Equation 7, with correction per Equation 15
		stepper.c0 = 0.676 * math.Sqrt(2.0/acceleration) * 1000000.0 // Equation 15
		stepper.acceleration = acceleration
		stepper.computeNewSpeed()
	}
}

func (stepper *AccelStepper) computeNewSpeed() {

	distanceTo := stepper.distanceToGo() // +ve is clockwise from curent location

	stepsToStop := int((stepper.speed * stepper.speed) / (2.0 * stepper.acceleration)) // Equation 16

	if distanceTo == 0 && stepsToStop <= 1 {
		// We are at the target and its time to stop
		stepper.stepInterval = 0
		stepper.speed = 0.0
		stepper.n = 0
		return
	}

	if distanceTo > 0 {
		// We are anticlockwise from the target
		// Need to go clockwise from here, maybe decelerate now
		if stepper.n > 0 {
			// Currently accelerating, need to decel now? Or maybe going the wrong way?
			if (stepsToStop >= distanceTo) || stepper.direction == DirectionCCW {
				stepper.n = float64(-stepsToStop) // Start deceleration
			}
		} else if stepper.n < 0 {
			// Currently decelerating, need to accel again?
			if (stepsToStop < distanceTo) && stepper.direction == DirectionCW {
				stepper.n = -stepper.n // Start accceleration
			}
		} else if distanceTo < 0 {
			// We are clockwise from the target
			// Need to go anticlockwise from here, maybe decelerate
			if stepper.n > 0 {
				// Currently accelerating, need to decel now? Or maybe going the wrong way?
				if (stepsToStop >= -distanceTo) || stepper.direction == DirectionCW {
					stepper.n = float64(-stepsToStop) // Start deceleration
				}
			} else if stepper.n < 0 {
				// Currently decelerating, need to accel again?
				if (stepsToStop < -distanceTo) && stepper.direction == DirectionCCW {
					stepper.n = -stepper.n // Start accceleration
				}
			}
		}
	}

	// Need to accelerate or decelerate
	if stepper.n == 0 {
		// First step from stopped
		stepper.cn = stepper.c0
		if distanceTo > 0 {
			stepper.direction = DirectionCW
		} else {
			stepper.direction = DirectionCCW
		}
	} else {
		// Subsequent step. Works for accel (n is +_ve) and decel (n is -ve).
		stepper.cn = stepper.cn - ((2.0 * stepper.cn) / ((4.0 * stepper.n) + 1)) // Equation 13
		stepper.cn = math.Max(stepper.cn, stepper.cmin)
	}
	stepper.n = stepper.n + 1
	stepper.stepInterval = int(stepper.cn)
	stepper.speed = 1000000.0 / stepper.cn
	if stepper.direction == DirectionCCW {
		stepper.speed = -stepper.speed
	}
}

//if 0
//        Serial.println(_speed);
//        Serial.println(_acceleration);
//        Serial.println(_cn);
//        Serial.println(_c0);
//        Serial.println(_n);
//        Serial.println(_stepInterval);
//        Serial.println(distanceTo);
//        Serial.println(stepsToStop);
//        Serial.println("-----");
//endif

func (stepper *AccelStepper) distanceToGo() int {
	return stepper.targetPos - stepper.currentPos
}
func (stepper *AccelStepper) targetPosition() int {
	return stepper.targetPos
}
func (stepper *AccelStepper) currentPosition() int {
	return stepper.currentPos
}

// Useful during initialisations or after initial positioning
// Sets speed to 0
func (stepper *AccelStepper) setCurrentPosition(position int) {
	stepper.targetPos = position
	stepper.currentPos = position
	stepper.n = 0
	stepper.stepInterval = 0
	stepper.speed = 0.0
}

func (stepper *AccelStepper) moveTo(absolute int) {
	if stepper.targetPos != absolute {
		stepper.targetPos = absolute
		stepper.computeNewSpeed()

		// compute new n?
	}
}
func (stepper *AccelStepper) move(relative int) {
	stepper.moveTo(stepper.currentPos + relative)
}

func (stepper *AccelStepper) setMinPulseWidth(minWidth int) {
	stepper.minPulseWidth = minWidth
}

//func (stepper *AccelStepper) setEnablePin(enablePin byte) {
//	stepper.enablePin = enablePin

//	// This happens after construction, so init pin now.
//	if stepper.enablePin != 0xff {
//		GPIO.setup(stepper.enablePin, GPIO.OUT)
//		GPIO.output(stepper.enablePin, GPIO.HIGH^stepper.enableInverted)
//	}
//}

func (stepper *AccelStepper) setPinsInvertedDirStep(directionInvert, stepInvert int, enableInvert bool) { //func (stepper *AccelStepper) setPinsInverted(stepper, directionInvert, stepInvert, enableInvert):
	stepper.pinInverted[0] = stepInvert
	stepper.pinInverted[1] = directionInvert
	stepper.enableInverted = enableInvert
}

//// Prevents power consumption on the outputs
//func (stepper *AccelStepper) disableOutputs() {
//	if stepper.interf == 0 {
//		return
//	}
//
//	stepper.setOutputPins(0) // Handles inversion automatically
//	if stepper.enablePin != 0xff {
//
//		GPIO.setup(stepper.enablePin, GPIO.OUT)
//		GPIO.output(stepper.enablePin, GPIO.LOW^stepper.enableInverted)
//	}
//}

func (stepper *AccelStepper) setPinsInverted(pin1Invert, pin2Invert, pin3Invert, pin4Invert int, enableInvert bool) {
	stepper.pinInverted[0] = pin1Invert
	stepper.pinInverted[1] = pin2Invert
	stepper.pinInverted[2] = pin3Invert
	stepper.pinInverted[3] = pin4Invert
	stepper.enableInverted = enableInvert
}

// Blocks until the target position is reached and stopped
func (stepper *AccelStepper) runToPosition() {
	for stepper.run() {
	}
}
func (stepper *AccelStepper) runSpeedToPosition() bool {
	if stepper.targetPos == stepper.currentPos {
		return false
	}
	if stepper.targetPos > stepper.currentPos {
		stepper.direction = DirectionCW
	} else {
		stepper.direction = DirectionCCW
	}
	return stepper.runSpeed()
}

// Blocks until the new target position is reached
func (stepper *AccelStepper) runToNewPosition(position int) {
	stepper.moveTo(position)
	stepper.runToPosition()
}

func (stepper *AccelStepper) stop() {
	if stepper.speed != 0.0 {
		stepsToStop := ((stepper.speed * stepper.speed) / (2.0 * stepper.acceleration)) + 1 // Equation 16 (+integer rounding)
		if stepper.speed > 0 {
			stepper.move(int(stepsToStop))
		} else {
			stepper.move(int(-stepsToStop))
		}
	}
}

func (stepper *AccelStepper) isRunning() bool {
	return !(stepper.speed == 0.0 && stepper.targetPos == stepper.currentPos)
}
