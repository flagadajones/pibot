package main

func main() {
	v := Visage{}
	v.Init(&CapSize{w: 100, h: 100})
	v.draw(&Cible{x: 0, y: 0, w: 212, h: 120})
}
