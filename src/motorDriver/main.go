package main

func main() {
	v := AccelStepper{}
	v.init(FULL4WIRE, 37,35, 33, 31, true)
	v.move(20)
	}
