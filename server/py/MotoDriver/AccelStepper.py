from enum import Enum
import RPi.GPIO as GPIO #pip install RPi.GPIO
import math
from datetime import datetime
import time
def constrain(val, min_val, max_val):
    return min(max_val, max(min_val, val))

class Direction(Enum):
    DIRECTION_CCW = 0 # Counter-Clockwise
    DIRECTION_CW = 1 # Clockwise

class MotorInterfaceType(Enum):
    FUNCTION  = 0 #< Use the functional interface, implementing your own driver functions (internal use only)
    DRIVER    = 1 #< Stepper Driver, 2 driver pins required
    FULL2WIRE = 2 #< 2 wire stepper, 2 motor pins required
    FULL3WIRE = 3 #< 3 wire stepper, such as HDD spindle, 3 motor pins required
    FULL4WIRE = 4 #< 4 wire full stepper, 4 motor pins required
    HALF3WIRE = 6 #< 3 wire half stepper, such as HDD spindle, 3 motor pins required
    HALF4WIRE = 8 #< 4 wire half stepper, 4 motor pins required


class AccelStepper(object):
    __pin = []
    __pinInverted=[]
    def __init__(self, interface,  pin1, pin2, pin3, pin4, enable):

        self.__interface = interface
        self.__currentPos = 0
        self.__targetPos = 0
        self.__speed = 0.0
        self.__maxSpeed = 1.0
        self.__acceleration = 0.0
        self.__sqrt_twoa = 1.0
        self.__stepInterval = 0
        self.__minPulseWidth = 1
        self.__enablePin = 0xff
        self.__lastStepTime = 0
        self.__pin.append(pin1)
        self.__pin.append(pin2)
        self.__pin.append(pin3)
        self.__pin.append(pin4)
        self.__enableInverted = False


        self.__forward=lambda *_, **_: None
        self.__backward=lambda *_, **_: None


        
        # NEW
        self.__n = 0
        self.__c0 = 0.0
        self.__cn = 0.0
        self.__cmin = 1.0
        self.__direction = Direction.DIRECTION_CCW

        for i in range(0, 4):
            self.__pinInverted[i] = 0
        if (enable):
            self.enableOutputs()
        # Some reasonable default
        self.setAcceleration(1)
    
    def setSpeed(self, speed):
        if (speed == self.__speed):
            return
        speed = constrain(speed, -self.__maxSpeed, self.__maxSpeed)
        if (speed == 0.0):
            self.__stepInterval = 0
        else:
            self.__stepInterval = math.fabs(1000000.0 / speed)
            self.__direction =  Direction.DIRECTION_CW if (speed > 0.0) else Direction.DIRECTION_CCW
        
        self.__speed = speed

    def speed(self):
        return self.__speed


    # Subclasses can override
    def step(self, step):

        if(self.__interface == MotorInterfaceType.FUNCTION):
            self.step0(step)
            return
        if(self.__interface == MotorInterfaceType.DRIVER):
            self.step1(step)
            return
        if(self.__interface == MotorInterfaceType.FULL2WIRE):
            self.step2(step)
            return
        if(self.__interface == MotorInterfaceType.FULL3WIRE):
            self.step3(step)
            return
        if(self.__interface == MotorInterfaceType.FULL4WIRE):
            self.step4(step)
            return
        if(self.__interface == MotorInterfaceType.HALF3WIRE):
            self.step6(step)
            return
        if(self.__interface == MotorInterfaceType.HALF4WIRE):
            self.step8(step)
            return

    # You might want to override this to implement eg serial output
    # bit 0 of the mask corresponds to _pin[0]
    # bit 1 of the mask corresponds to _pin[1]
    # ....
    def setOutputPins(self, mask):
        numpins = 2
        if (self.__interface == MotorInterfaceType.FULL4WIRE or self.__interface == MotorInterfaceType.HALF4WIRE):
            numpins = 4
        elif (self.__interface == MotorInterfaceType.FULL3WIRE or self.__interface == MotorInterfaceType.HALF3WIRE):
            numpins = 3
        
        for i in range ( 0, numpins):
            GPIO.output(self.__pin[i], (GPIO.HIGH ^ self.__pinInverted[i]) if  (mask & (1 << i)) else (GPIO.LOW ^ self.__pinInverted[i])) # 	digitalWrite(_pin[i], (mask & (1 << i)) ? (HIGH ^ _pinInverted[i]) : (LOW ^ _pinInverted[i]));

    # 0 pin step function (ie for functional usage)
    def step0(self, step):
        #(void)(step); // Unused
        if (self.__speed > 0):
            self.__forward()
        else:
            self.__backward()

    # 1 pin step function (ie for stepper drivers)
    # This is passed the current step number (0 to 7)
    # Subclasses can override
    def step1(self, step):
        # (void)(step); // Unused
        # _pin[0] is step, _pin[1] is direction
        self.setOutputPins(0b10 if self.__direction else 0b00) # Set direction first else get rogue pulses
        self.setOutputPins(0b11 if self.__direction else 0b01) # step HIGH
        # Caution 200ns setup time 
        # Delay the minimum allowed pulse width
        time.sleep(self.__minPulseWidth *0.000001) # delayMicroseconds(self.__minPulseWidth)
        self.setOutputPins(0b10 if self.__direction else 0b00) # step LOW

    # 2 pin step function
    # This is passed the current step number (0 to 7)
    # Subclasses can override
    def step2(self, step):
        value=step & 0x3
        if( value==0): # 01 
            self.setOutputPins(0b10)
            return
        if( value==1): # 11 
            self.setOutputPins(0b11)
            return
        if( value==2): # 10 
            self.setOutputPins(0b01)
            return
        if( value==3): # 00 
            self.setOutputPins(0b00)
            return

    # 3 pin step function
    # This is passed the current step number (0 to 7)
    # Subclasses can override
    def step3(self, step):
        value = step % 3
        if(value ==0 ): # 100
            self.setOutputPins(0b100)
            return
        if(value ==1 ): # 001
            self.setOutputPins(0b001)
            return
        if(value ==2 ): # 010
            self.setOutputPins(0b010)
            return

    # 4 pin step function for half stepper
    # This is passed the current step number (0 to 7)
    # Subclasses can override
    def step4(self,step):
        value = step & 0x3
        if(value ==0 ): # 1010
            self.setOutputPins(0b0101)
            return
        if(value ==1 ): # 0110
            self.setOutputPins(0b0110)
            return
        if(value ==2 ): # 0101
            self.setOutputPins(0b1010)
            return
        if(value ==3 ): # 1001
            self.setOutputPins(0b1001)
            return

    # 3 pin half step function
    # This is passed the current step number (0 to 7)
    # Subclasses can override
    def step6(self, step):
        value = step % 6
        if(value ==0 ): # 100
            self.setOutputPins(0b100)
            return
        if(value ==1 ): # 101
            self.setOutputPins(0b101)
            return
        if(value ==2 ): # 001
            self.setOutputPins(0b001)
            return
        if(value ==3 ): # 011
            self.setOutputPins(0b011)
            return
        if(value ==4 ): # 010
            self.setOutputPins(0b010)
            return
        if(value ==5 ): # 110
            self.setOutputPins(0b110)
            return

    # 4 pin half step function
    # This is passed the current step number (0 to 7)
    # Subclasses can override
    def step8(self, step):
        value = step & 0x7
        if(value ==0 ): # 1000
            self.setOutputPins(0b0001)
            return
        if(value ==1 ): # 1010
            self.setOutputPins(0b0101)
            return
        if(value ==2 ): # 0010
            self.setOutputPins(0b0100)
            return
        if(value ==3 ): # 0110
            self.setOutputPins(0b0110)
            return
        if(value ==4 ): # 0100
            self.setOutputPins(0b0010)
            return
        if(value ==5 ): # 0101
            self.setOutputPins(0b1010)
            return
        if(value ==6 ): # 0001
            self.setOutputPins(0b1000)
            return
        if(value ==7 ): # 1001
            self.setOutputPins(0b1001)
            return
    



    # Implements steps according to the current step interval
    # You must call this at least once per step
    # returns true if a step occurred
    def runSpeed(self):
        # Dont do anything unless we actually have a step interval
        if (not self.__stepInterval):
            return False

        timer = datetime.now().microseconds()   
        if (timer - self.__lastStepTime >= self.__stepInterval):
            if (self.__direction == Direction.DIRECTION_CW):
                # Clockwise
                self.__currentPos += 1
            else:
                # Anticlockwise  
                self.__currentPos -= 1
            
            self.step(self.__currentPos)

            _lastStepTime = timer # Caution: does not account for costs in step()

            return True
        else:
            return False
    
    # Run the motor to implement speed and acceleration in order to proceed to the target position
    # You must call this at least once per step, preferably in your main loop
    # If the motor is in the desired position, the cost is very small
    # returns true if the motor is still running to the target position.
    def run(self):
        if (self.runSpeed()):
            self.computeNewSpeed()
        return self.__speed != 0.0 or self.distanceToGo() != 0


    def setMaxSpeed(self, speed):
        if (speed < 0.0):
            speed = -speed
        if (self.__maxSpeed != speed):
            self.__maxSpeed = speed
            self.__cmin = 1000000.0 / speed
            # Recompute _n from current speed and adjust speed if accelerating or cruising
            if (self.__n > 0):
                self.__n = ((self.__speed * self.__speed) / (2.0 * self.__acceleration)) # Equation 16
                self.computeNewSpeed()

    def maxSpeed(self):
        return self.__maxSpeed



    def enableOutputs(self):

        if (not self.__interface): 
            return

        GPIO.setup(self.__pin[0], GPIO.OUT)
        GPIO.setup(self.__pin[1], GPIO.OUT)
        if (self.__interface == MotorInterfaceType.FULL4WIRE or self.__interface == MotorInterfaceType.HALF4WIRE):
            GPIO.setup(self.__pin[2], GPIO.OUT)
            GPIO.setup(self.__pin[3], GPIO.OUT)
        
        elif (self.__interface == MotorInterfaceType.FULL3WIRE or self.__interface == MotorInterfaceType.HALF3WIRE):
            GPIO.setup(self.__pin[2], GPIO.OUT)
      
        if (self.__enablePin is not 0xff):
            GPIO.setup(self.__enablePin, GPIO.OUT)
            GPIO.output(self.__enablePin, GPIO.HIGH ^ self.__enableInverted) #   digitalWrite(_enablePin, GPIO.HIGH ^ self.__enableInverted);
        

    def setAcceleration(self, acceleration):

        if (acceleration == 0.0):
            return
        if (acceleration < 0.0):
            acceleration = -acceleration
        if (self.__acceleration is not acceleration):
            # Recompute _n per Equation 17
            self.__n = self.__n * (self.__acceleration / acceleration)
            # New c0 per Equation 7, with correction per Equation 15
            self.__c0 = 0.676 * math.sqrt(2.0 / acceleration) * 1000000.0 # Equation 15
            self.__acceleration = acceleration
            self.computeNewSpeed()
        
    def computeNewSpeed(self):

        distanceTo = self.distanceToGo() # +ve is clockwise from curent location

        stepsToStop = ((self.__speed * self.__speed) / (2.0 * self.__acceleration)) # Equation 16

        if (distanceTo == 0 and stepsToStop <= 1):
            # We are at the target and its time to stop
            self.__stepInterval = 0
            self.__speed = 0.0
            self.__n = 0
            return
            

        if (distanceTo > 0):
            # We are anticlockwise from the target
            # Need to go clockwise from here, maybe decelerate now
            if (self.__n > 0):
                # Currently accelerating, need to decel now? Or maybe going the wrong way?
                if ((stepsToStop >= distanceTo) or self.__direction == Direction.DIRECTION_CCW):
                    self.__n = -stepsToStop # Start deceleration
            elif (self.__n < 0):
                # Currently decelerating, need to accel again?
                if ((stepsToStop < distanceTo) and self.__direction == Direction.DIRECTION_CW):
                    self.__n = -self.__n # Start accceleration
        elif (distanceTo < 0):
            # We are clockwise from the target
            # Need to go anticlockwise from here, maybe decelerate
            if (self.__n > 0):
                # Currently accelerating, need to decel now? Or maybe going the wrong way?
                if ((stepsToStop >= -distanceTo) or self.__direction == Direction.DIRECTION_CW):
                    self.__n = -stepsToStop # Start deceleration
            elif (self.__n < 0):
                # Currently decelerating, need to accel again?
                if ((stepsToStop < -distanceTo) and self.__direction == Direction.DIRECTION_CCW):
                    self.__n = -self.__n# Start accceleration

        # Need to accelerate or decelerate
        if (self.__n == 0):
            # First step from stopped
            self.__cn = self.__c0
            self.__direction = Direction.DIRECTION_CW if (distanceTo > 0) else Direction.DIRECTION_CCW
        else:
            # Subsequent step. Works for accel (n is +_ve) and decel (n is -ve).
            self.__cn = self.__cn - ((2.0 * self.__cn) / ((4.0 * self.__n) + 1)) # Equation 13
            self.__cn = max(self.__cn, self.__cmin)

        self.__n=self.__n+1
        self.__stepInterval = self.__cn
        self.__speed = 1000000.0 / self.__cn
        if (self.__direction == Direction.DIRECTION_CCW):
            self.__speed = -self.__speed

    #if 0
#        Serial.println(_speed);
#        Serial.println(_acceleration);
#        Serial.println(_cn);
#        Serial.println(_c0);
#        Serial.println(_n);
#        Serial.println(_stepInterval);
#        Serial.println(distanceTo);
#        Serial.println(stepsToStop);
#        Serial.println("-----");
    #endif

    def distanceToGo(self):
       return self.__targetPos - self.__currentPos

    def targetPosition(self):
        return self.__targetPos

    def currentPosition(self):
        return self.__currentPos


# Useful during initialisations or after initial positioning
# Sets speed to 0
    def setCurrentPosition(self, position):
        self.__targetPos = self.__currentPos = position
        self.__n = 0
        self.__stepInterval = 0
        self.__speed = 0.0

    def moveTo(self, absolute):
        if (self.__targetPos != absolute):
            self.__targetPos = absolute
            self.computeNewSpeed()
            # compute new n?
    
    def move(self, relative):
        self.moveTo(self.__currentPos + relative)


    def setMinPulseWidth(self, minWidth):
        self.__minPulseWidth = minWidth

    def setEnablePin(self, enablePin):
        self.__enablePin = enablePin

        # This happens after construction, so init pin now.
        if (self.__enablePin != 0xff):
            GPIO.setup(self.__enablePin, GPIO.OUT)
            GPIO.output(self.__enablePin, GPIO.HIGH ^ self.__enableInverted)


    def setPinsInvertedDirStep(self, directionInvert, stepInvert, enableInvert): #def setPinsInverted(self, directionInvert, stepInvert, enableInvert):
        self.__pinInverted[0] = stepInvert
        self.__pinInverted[1] = directionInvert
        self.__enableInverted = enableInvert

    # Prevents power consumption on the outputs
    def disableOutputs(self):
        if (not  self.__interface): 
            return

        self.setOutputPins(0) # Handles inversion automatically
        if (self.__enablePin != 0xff):
            GPIO.setup(self.__enablePin, GPIO.OUT)
            GPIO.output(self.__enablePin, GPIO.LOW ^ self.__enableInverted)


    def setPinsInverted(self, pin1Invert, pin2Invert, pin3Invert, pin4Invert, enableInvert):
        self.__pinInverted[0] = pin1Invert
        self.__pinInverted[1] = pin2Invert
        self.__pinInverted[2] = pin3Invert
        self.__pinInverted[3] = pin4Invert
        self.__enableInverted = enableInvert

    # Blocks until the target position is reached and stopped
    def runToPosition(self):
        while (self.run()):
            pass

    def runSpeedToPosition(self):
        if (self.__targetPos == self.__currentPos):
            return False
        if (self.__targetPos >self.__currentPos):
            self.__direction = Direction.DIRECTION_CW
        else:
            self.__direction = Direction.DIRECTION_CCW
        return self.runSpeed()


    # Blocks until the new target position is reached
    def runToNewPosition(self, position):
        self.moveTo(position)
        self.runToPosition()

    def stop(self):
        if (self.__speed != 0.0):
            stepsToStop = ((self.__speed * self.__speed) / (2.0 * self.__acceleration)) + 1 # Equation 16 (+integer rounding)
            if (self.__speed > 0):
                self.move(stepsToStop)
            else:
                self.move(-stepsToStop)

    def isRunning(self):
        return not(self.__speed == 0.0 and self.__targetPos == self.__currentPos)
