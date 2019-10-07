package main

import (
	"log"
	"math"

	"github.com/flagadajones/pibot/visage"
)

type Point struct {
	x float64
	y float64
	z float64
}

func main() {
	//5cm  ==> 200px

	//10cm ==> 100px

	//30cm ==> 33px
	visageA30cm := visage.Cible{X: 1, Y: 1, W: 477, H: 477}
	//visageA10cm := visage.Cible{X: 1, Y: 1, W: 100, H: 100}
	cible := visage.Cible{X: 1, Y: 1, W: 30, H: 30}
	var pointM Point
	rayon := float64(10)
	m := Point{x: float64(cible.X + cible.W/2), y: float64(cible.Y + cible.H/2), z: 10 / (float64(cible.W) / float64(visageA10cm.W))}
	om := math.Sqrt(float64(m.x*m.x + m.y*m.y + m.z*m.z))
	mh := float64(m.z)
	sinZigm := mh / om
	pointM.z = rayon * sinZigm

	oi := float64(m.x)
	oh := math.Sqrt(float64(m.x*m.x) + float64(m.y*m.y))
	cosTheta := oi / oh
	cosZigm := oh / om
	pointM.x = rayon * cosTheta * cosZigm

	ih := float64(m.y)
	sinTheta := ih / oh
	pointM.y = rayon * cosZigm * sinTheta

	log.Println(pointM)

	//ratio := float32(visageA50cm.W) / float32(visageA10cm.W)
	//log.Println(10 / ratio)
	//log.Println(visageA10cm, visageA50cm, ratio)

}
