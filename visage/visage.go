package visage

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Point struct {
	x float64
	y float64
	z float64
}
type Cible struct {
	X, Y, W, H int
}

func (v *Visage) calcul(cible *Cible, o Point) Point {
	var pointM Point
	cmPerPx := 4.2 / 213
	rayon := float64(69+69) / 2 * cmPerPx
	visageA30cm := Cible{X: 1, Y: 1, W: 477, H: 477}
	x := cible.X
	y := cible.Y
	w := cible.W
	hh := cible.H
	println("==== ", x, y, w, hh)
	println("== / ", v.size.W, v.capWidth, v.size.H, v.capHeight)
	v.cibleX = (x + w/2) * v.size.W / v.capWidth
	v.cibleY = (y + hh/2) * v.size.H / v.capHeight
	v.cibleW = cible.W * v.size.W / v.capWidth
	v.cibleH = cible.H * v.size.H / v.capHeight
	println("===+ ", v.cibleX, v.cibleY, v.cibleW, v.cibleH)
	dist := float64(5)
	//rayon := float64(v.eyeHeight) * v.cmPerPx
	//log.Print("O", o)

	//log.Print("O", o)
	//	v.cibleX = 424
	//	v.cibleY = 240
	//log.Print("cibleXY")
	//log.Print(float64(v.cibleX)*v.cmPerPx, float64(v.cibleY)*v.cmPerPx)
	m := Point{x: (float64(v.cibleX) - o.x) * cmPerPx, y: (float64(v.cibleY) - o.y) * cmPerPx, z: dist / (float64(v.cibleW) / float64(visageA30cm.W))}
	o.x = o.x * cmPerPx
	o.y = o.y * cmPerPx
	o.z = o.z * cmPerPx
	o.x = 0
	o.y = 0
	o.z = 0
	//log.Print("m", m)

	detX := 1.0
	if m.x < 0 {
		detX = -1.0
	}
	detY := 1.0
	if m.y < 0 {
		detY = -1.0

	}
	om := math.Sqrt(float64((m.x-o.x)*(m.x-o.x) + (m.y-o.y)*(m.y-o.y) + (m.z-o.z)*(m.z-o.z)))
	//log.Print("om", om)
	//√ ((x_B - x_A)^2 + (y_B - y_A)^2 + (z_B - z_A)^2
	h := Point{x: m.x, y: m.y, z: o.z}
	//log.Print("h", h)

	mh := math.Sqrt(float64((h.x-m.x)*(h.x-m.x) + (h.y-m.y)*(h.y-m.y) + (h.z-m.z)*(h.z-m.z)))
	//log.Print("mh", mh)

	sinZigm := mh / om
	//log.Print("sinZ", sinZigm)
	pointM.z = rayon * sinZigm / cmPerPx
	i := Point{x: h.x, y: o.y, z: o.z}
	//log.Print("i", i)

	oi := math.Sqrt(float64((i.x-o.x)*(i.x-o.x) + (i.y-o.y)*(i.y-o.y) + (i.z-o.z)*(i.z-o.z)))

	//log.Print("oi", oi)
	oh := math.Sqrt(float64((h.x-o.x)*(h.x-o.x) + (h.y-o.y)*(h.y-o.y) + (h.z-o.z)*(h.z-o.z)))
	//log.Print("oh", oh)

	cosTheta := detX * oi / oh
	//log.Print("cosT", cosTheta)
	cosZigm := oh / om
	//log.Print("cosZ", cosTheta)
	pointM.x = rayon * cosTheta * cosZigm / cmPerPx

	ih := math.Sqrt(float64((h.x-i.x)*(h.x-i.x) + (h.y-i.y)*(h.y-i.y) + (h.z-i.z)*(h.z-i.z)))
	sinTheta := detY * ih / oh
	//log.Print("sinT", sinTheta)
	pointM.y = rayon * cosZigm * sinTheta / cmPerPx
	//log.Print("pM", pointM)
	return pointM
}
func (v *Visage) Draw(cible *Cible) {
	v.cible = cible

}

type CapSize struct {
	W, H int
}

func (v *Visage) Init(capSize *CapSize) {
	rl.InitWindow(656, 416, "pibot")
	v.cible = &Cible{X: 320, Y: 200, W: 50, H: 50}
	v.capWidth = capSize.W
	v.capHeight = capSize.H
	v.size = CapSize{W: 656, H: 416}
}

func (v *Visage) Blink() {
}

type Visage struct {
	cible                          *Cible
	capWidth, capHeight            int
	size                           CapSize
	cibleX, cibleY, cibleW, cibleH int
}

func (v *Visage) Run() {
	//rl.SetConfigFlags(rl.FlagFullscreenMode)

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)

	textureV := rl.LoadTexture("/Users/grua341/go/src/github.com/flagadajones/pibot/visage/visage.png")   // Texture loading
	textureF := rl.LoadTexture("/Users/grua341/go/src/github.com/flagadajones/pibot/visage/fondoeil.png") // Texture loading
	textureI := rl.LoadTexture("/Users/grua341/go/src/github.com/flagadajones/pibot/visage/newIris.png")  // Texture loading

	rl.SetTargetFPS(30)

	X1 := float64(150 - textureI.Width/2)
	Y1 := float64(155 - textureI.Height/2)
	X2 := float64(509 - textureI.Width/2)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		//		cible := Cible{X: 0, Y: 0, W: 656, H: 416}

		//cible := Cible{X: 300, Y: 200, W: 50, H: 50}

		//	cible := Cible{X: 600, Y: 200, W: 50, H: 50}
		//cible := Cible{X: 600, Y: 400, W: 50, H: 50}

		//cible := Cible{X: 0, Y: 400, W: 50, H: 50}

		p1 := v.calcul(v.cible, Point{x: X1, y: Y1, z: 0})
		p2 := v.calcul(v.cible, Point{x: X2, y: Y1, z: 0})

		rl.DrawTexture(textureF, 0, 0, rl.RayWhite)
		rl.DrawTexture(textureI, int32(p1.x+X1), int32(p1.y+Y1), rl.RayWhite)
		rl.DrawTexture(textureI, int32(p2.x+X2), int32(p2.y+Y1), rl.RayWhite)
		rl.DrawTexture(textureV, 0, 0, rl.RayWhite)
		rl.DrawCircle(int32(v.cibleX), int32(v.cibleY), 10, rl.Red)
		rl.DrawRectangleLines(int32(v.cibleX-v.cibleW/2), int32(v.cibleY-v.cibleH/2), int32(v.cibleW), int32(v.cibleH), rl.Blue)
		//println(v.cible.X, v.cible.Y, v.cible.W, v.cible.H)
		//rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}
	rl.UnloadTexture(textureV) // Texture unloading
	rl.UnloadTexture(textureF)
	rl.UnloadTexture(textureI)
	rl.CloseWindow()

}
