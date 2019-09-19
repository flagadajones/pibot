from AccelStepper import AccelStepper,MotorInterfaceType
import time
import math
 
BASE_STEPS_PER_REV=200
MAX_SPEED=9600
MAX_ACCEL=500
MIN_PULSE =20

class OmniWheelDriver:

    def __init__(self,wheel_diam,base_diam, max_speed=MAX_SPEED,max_accel=MAX_ACCEL):
        self.stepper0=AccelStepper(MotorInterfaceType.FULL4WIRE, 1, 2,3,4,True)
        self.stepper1=AccelStepper(MotorInterfaceType.FULL4WIRE, 1, 2,3,4,True)
        self.stepper2=AccelStepper(MotorInterfaceType.FULL4WIRE, 1, 2,3,4,True)
        self.stepper3=AccelStepper(MotorInterfaceType.FULL4WIRE, 1, 2,3,4,True)
        self._index=0

        self._step_=0b001
        self._sleep_=True
        self._enable_=False

        self._max_speed_=max_speed
        self._max_accel_=max_accel
        self._wheel_diameter_=wheel_diam
        self._base_diameter_=base_diam
        
        self.stepping_table=[]
        self.steps_per_rev=[]
        self.stepping_table[0] = "Full Step"
        self.steps_per_rev[0] = BASE_STEPS_PER_REV
        self.stepping_table[1] = "Half Step"
        self.steps_per_rev[1] = 2 * BASE_STEPS_PER_REV
        self.stepping_table[2] = "Quarter Step"
        self.steps_per_rev[2] = 4 * BASE_STEPS_PER_REV
        self.stepping_table[3] = "Eighth Step"
        self.steps_per_rev[3] = 8 * BASE_STEPS_PER_REV
        self.stepping_table[4] = "UNDEFINED Using Full Step"
        self.steps_per_rev[4] = 0
        self.stepping_table[5] = "UNDEFINED Using Full Step"
        self.steps_per_rev[5] = 0
        self.stepping_table[6] = "UNDEFINED Using Full Step"
        self.steps_per_rev[6] = 0
        self.stepping_table[7] = "Sisteenth Step"
        self.steps_per_rev[7] = 16 * BASE_STEPS_PER_REV


      #  pinMode(ENABLE, OUTPUT)
      #  pinMode(SLEEP, OUTPUT)
      #  pinMode(RESET, OUTPUT)
      #  pinMode(MS3, OUTPUT)
      #  pinMode(MS2, OUTPUT)
      #  pinMode(MS1, OUTPUT)
        
        self.reset_driver()
        self.set_step(self._step_)
        self.set_max_speed(self._max_speed_)
        self.set_max_accel(self._max_accel_)
        self.set_min_pulse_width(MIN_PULSE)
        self.zero_location()
        self.disable_steppers()
        self.sleep()

    def display_status(self):
        print("SLEEP: ")
        print(self._sleep_)
        print("ENABLE: ")
        print(self._enable_)
        print("OMNI WHEEL DIAMETER: ")
        print(self._wheel_diameter_)
        print("cm")
        print("BASE DIAMETER: ")
        print(self._base_diameter_)
        print("cm")
        print("step = ")
        print(self.stepping_table[self._step_])
        print("steps per rev = ")
        print(self.steps_per_rev[self._step_])
        print("max speed = ")
        print(self._max_speed_)

    
    def set_base_diameter(self, diameter):
        self._base_diameter_ = diameter
        print("setting base diameter to ")
        print(self._base_diameter_)
        print("cm")


    def set_wheel_diameter(self, diameter):
        self._wheel_diameter_ = diameter
        print("setting omni wheel diameter to ")
        print(self._wheel_diameter_)
        print("cm")

    def set_step(self, step):
        self._step_ = step
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
        print("step = ")
        print(self.stepping_table[self._step_])
        print("steps per rev = ")
        print(self.steps_per_rev[self._step_])

    def set_max_speed(self, speed):
        self._max_speed_ = speed
        self.stepper0.setMaxSpeed(speed)
        self.stepper1.setMaxSpeed(speed)
        self.stepper2.setMaxSpeed(speed)
        self.stepper3.setMaxSpeed(speed)
        
        print("max speed set to ")
        print(speed) 




    def set_max_accel(self, accel):
        self._max_accel_ = accel
        self.stepper0.setAcceleration(accel)
        self.stepper1.setAcceleration(accel)
        self.stepper2.setAcceleration(accel)
        self.stepper3.setAcceleration(accel)
        
        print("max accel set to ")
        print(accel)


    def set_min_pulse_width(self, pulse_width):
        self.stepper0.setMinPulseWidth(pulse_width)
        self.stepper1.setMinPulseWidth(pulse_width)
        self.stepper2.setMinPulseWidth(pulse_width)
        self.stepper3.setMinPulseWidth(pulse_width)



    def enable_steppers(self):
        #digitalWrite(ENABLE, LOW)
        self._enable_ = True


    def disable_steppers(self):
        #digitalWrite(ENABLE, HIGH);
        self._enable_ = False


    def sleep(self): 
        #digitalWrite(SLEEP, LOW)
        self._sleep_ = True


    def wake(self):
        #digitalWrite(SLEEP, HIGH)
        self._sleep_ = False

    def zero_location(self):
        self.stepper0.setCurrentPosition(0)
        self.stepper1.setCurrentPosition(0)
        self.stepper2.setCurrentPosition(0)
        self.stepper3.setCurrentPosition(0)

    def reset_driver(self):
        #digitalWrite(RESET, LOW)
        time.sleep(0.100)
        #digitalWrite(RESET, HIGH):

    def hard_stop_all(self):
        self.stepper0.setCurrentPosition(self.stepper0.targetPosition())
        self.stepper1.setCurrentPosition(self.stepper1.targetPosition())
        self.stepper2.setCurrentPosition(self.stepper2.targetPosition())
        self.stepper3.setCurrentPosition(self.stepper3.targetPosition())

    def run(self):
        self.stepper0.run()
        self.stepper1.run()
        self.stepper2.run()
        self.stepper3.run()

    def stop_all(self):
        self.stepper0.stop()
        self.stepper1.stop()
        self.stepper2.stop()
        self.stepper3.stop()

    def rotate(self, deg):
        self.set_max_accel(self._max_accel_)
        base_circumference = math.pi * self._base_diameter_
        wheel_circumference = math.pi * self._wheel_diameter_
        cm_per_step = wheel_circumference / self.steps_per_rev[self._step_]
        steps = ((deg / 360.0) * base_circumference) / cm_per_step
        
        self.stepper0.move(steps)
        self.stepper1.move(steps)
        self.stepper2.move(steps)
        self.stepper3.move(steps)

    def moveto(self, r, theta):
        theta_rad = math.pi / 180.0 * theta
        x = r * math.cos(theta_rad)
        y = r * math.sin(theta_rad)
        wheel_circumference = math.pi * self._wheel_diameter_
        cm_per_step = wheel_circumference / self.steps_per_rev[self._step_]
        x_steps = x / cm_per_step
        y_steps = y / cm_per_step

        if (math.fabs(x_steps) > math.fabs(y_steps)):
            slow_percent = math.fabs(y_steps / x_steps) 
            slow_speed = slow_percent * self._max_accel_
            self.stepper0.setAcceleration(self._max_accel_)
            self.stepper1.setAcceleration(self._max_accel_)
            self.stepper2.setAcceleration(slow_speed)
            self.stepper3.setAcceleration(slow_speed) 
        elif (math.fabs(x_steps) < math.fabs(y_steps)):
            slow_percent = math.fabs(x_steps / y_steps) 
            slow_speed = slow_percent * self._max_accel_
            self.stepper0.setAcceleration(slow_speed)
            self.stepper1.setAcceleration(slow_speed)
            self.stepper2.setAcceleration(self._max_accel_)
            self.stepper3.setAcceleration(self._max_accel_)
        else:
            self.set_max_accel(self._max_accel_)
        
        self.stepper0.move(x_steps)
        self.stepper1.move(-x_steps)
        self.stepper2.move(y_steps)
        self.stepper3.move(-y_steps)

