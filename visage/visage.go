package visage

import (
	//"image/color"
	"image/png"
	"math"
	"os"
	"path"
	"runtime"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	//s"golang.org/x/exp/shiny/materialdesign/colornames"
)

var visage *Visage

// Run ...
func (v *Visage) Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Eye",
		Bounds: pixel.R(0, 0, float64(visage.size.W), float64(visage.size.H)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	//	win.SetMonitor(pixelgl.PrimaryMonitor())
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Black)
	_, currentfile, _, _ := runtime.Caller(0)
	iris := path.Join(path.Dir(currentfile), "iris.png")

	f, err := os.Open(iris)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	pd := pixel.PictureDataFromImage(img)
	visage.pupil = pixel.NewSprite(pd, pd.Bounds())
	visage.pd = pd
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	//log.Print("radiu")
	//log.Print(visage.eyeRadius)
	mat = mat.Scaled(win.Bounds().Center(), float64(visage.eyeRadius)/pd.Bounds().W())
	visage.pupil.Draw(win, mat)
	//	v.pupil = pygame.image.load(dir_path+'/iris.png')
	//		v.pupil = v.pupil.convert_alpha()
	//			v.pupil = pygame.transform.scale(
	//				v.pupil, [v.eyeRadius, v.eyeRadius])
	//			v.rad = v.pupil.get_width()/2
	visage.rad = int(visage.eyeRadius)

	//visage.draw(&Cible{x: 0, y: 0, w: 212, h: 120})
	for !win.Closed() {

		visage.composite(win)
		win.Update()
	}
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type Cible struct {
	X, Y, W, H int
}

type CapSize struct {
	W, H int
}

// Draw ...
func (v *Visage) Draw(cible *Cible) {
	v.calculCible(cible)
	v.moveEyeLeft()
	v.moveEyeRight()

}

type Visage struct {
	capWidth, capHeight int
	size                CapSize
	maxVt, ratioVt      int
	startx, starty      int
	eyeRadius           float64
	eyeRatio            float64
	eyeWidth, eyeHeight float64
	maxBlink            int
	cible               *Cible

	//FIXME
	lastrect                          *CapSize
	rad                               int
	pupil                             *pixel.Sprite
	pd                                *pixel.PictureData
	eyeMoveRadius, eyeMoveRadiusWidth int
	cibleX, cibleY                    int
	leftEllipseX, leftEllipseY        float64
	leftEyeX, leftEyeY                float64
	leftX, leftY                      float64
	leftTx, leftTy, leftVx, leftVy    int

	rightEllipseX, rightEllipseY         float64
	rightEyeX, rightEyeY, rightX, rightY float64
	rightTx, rightTy, rightVx, rightVy   int

	xxx, xxS int
	blink    bool

	X1, X2, Y1, Y2 float64
	cmPerPx        float64
}

func (v *Visage) Blink() {
	v.blink = true
}
func (v *Visage) Init(capSize *CapSize) {
	visage = v
	v.capWidth = capSize.W
	v.capHeight = capSize.H
	v.size = CapSize{W: 848, H: 480}
	//v.size = CapSize{W: 213, H: 120}
	//FIXME revoir taille pour calcul en CM......pour que ca marche tout le temps
	//v.cmPerPx = 4.2 / float64(v.size.W)
	v.cmPerPx = 4.2 / 213

	v.maxVt = 50
	v.ratioVt = 2
	v.blink = true
	//	v.lastrect = nil

	v.starty = int(v.size.H / 2)
	v.startx = int(v.size.W / 2)
	v.cibleX = v.startx
	v.cibleY = v.starty

	v.eyeRadius = float64(v.size.H) / 4
	v.eyeRatio = 3.0 / 5.0
	v.eyeWidth = v.eyeRadius
	v.eyeHeight = float64(v.eyeRadius) * v.eyeRatio
	//log.Print(v.starty)
	//log.Print(v.eyeRadius)
	//log.Print(v.eyeRatio)
	v.maxBlink = int(v.eyeHeight) * 2

	////	v.pupil = pygame.image.load(dir_path+'/iris.png')
	////		v.pupil = v.pupil.convert_alpha()
	////			v.pupil = pygame.transform.scale(
	////				v.pupil, [v.eyeRadius, v.eyeRadius])
	////			v.rad = v.pupil.get_width()/2

	v.rad = int(v.eyeRadius)
	v.eyeMoveRadius = int(v.eyeHeight * 2 / 3)
	v.eyeMoveRadiusWidth = int(v.eyeWidth * 2 / 3)

	v.leftEllipseX = float64(v.startx - v.startx/2) //v.startx - v.startx/2- v.eyeRadius
	v.leftEllipseY = float64(v.starty)              //	v.leftEllipseY = int(float64(v.starty) - float64(v.eyeRadius)*v.eyeRatio)
	v.leftEyeX = v.leftEllipseX
	v.leftEyeY = v.leftEllipseY
	v.leftX = v.leftEyeX
	v.leftY = v.leftEyeY
	v.leftTx = int(v.leftX)
	v.leftTy = int(v.leftY)
	v.leftVx = 0
	v.leftVy = 0

	v.rightEllipseX = float64(v.startx + v.startx/2) //v.startx + v.startx/2 - v.eyeRadius
	v.rightEllipseY = float64(v.starty)              //	v.rightEllipseY = int(float64(v.starty) - float64(v.eyeRadius)*v.eyeRatio)
	v.rightEyeX = v.rightEllipseX
	v.rightEyeY = v.rightEllipseY
	v.rightX = v.rightEyeX
	v.rightY = v.rightEyeY
	v.rightTx = int(v.rightX)
	v.rightTy = int(v.rightY)
	v.rightVx = 0
	v.rightVy = 0

	v.xxx = v.maxBlink
	v.xxS = 1
	//	v.blink = false

}

/*func (v *Visage) run() {	pixelgl.Run(run)

}
*/
func sign(value int) float64 {
	if value > 0 {
		return 1.0
	}
	if value < 0 {
		return -1.0
	}
	return 0.0
}

func (v *Visage) calculRayon(start, cible, eye CapSize) (float64, float64) {
	Xa := start.W - eye.W
	Ya := start.H - eye.H
	Xb := cible.W - eye.W
	Yb := cible.H - eye.H

	Na := math.Sqrt(float64(Xa*Xa + Ya*Ya))
	Nb := math.Sqrt(float64(Xb*Xb + Yb*Yb))
	C := float64(Xa*Xb+Ya*Yb) / (Na * Nb)
	S := (Xa*Yb - Ya*Xb)
	angle := sign(S) * math.Acos(C)
	return angle, Nb
}

type Point struct {
	x float64
	y float64
	z float64
}

/* func (v *Visage) calcul(cible *Cible, point Point) Point {
	var pointM Point
	//rayon := float64(v.eyeWidth+v.eyeHeight) / 2
	visageA30cm := Cible{X: 1, Y: 1, W: 477, H: 477}

	dist := float64(30)
	rayon := float64(v.eyeHeight)
	m := Point{x: (float64(v.cibleX) - point.x) * v.cmPerPx, y: (float64(v.cibleY) - point.y) * v.cmPerPx, z: dist / (float64(cible.W) / float64(visageA30cm.W))}
	log.Println("m", m)
	om := math.Sqrt(float64(m.x*m.x + m.y*m.y + m.z*m.z))
	mh := float64(m.z)
	sinZigm := mh / om
	pointM.z = rayon * sinZigm / v.cmPerPx

	oi := float64(m.x)
	oh := math.Sqrt(float64(m.x*m.x) + float64(m.y*m.y))
	cosTheta := oi / oh
	cosZigm := oh / om
	pointM.x = rayon * cosTheta * cosZigm / v.cmPerPx

	ih := float64(m.y)
	sinTheta := ih / oh
	pointM.y = rayon * cosZigm * sinTheta / v.cmPerPx
	log.Println("pM", pointM)
	return pointM
} */

func (v *Visage) calcul(cible *Cible, o Point) Point {
	var pointM Point
	rayon := float64(v.eyeWidth+v.eyeHeight) / 2 * v.cmPerPx
	visageA30cm := Cible{X: 1, Y: 1, W: 477, H: 477}

	dist := float64(5)
	//rayon := float64(v.eyeHeight) * v.cmPerPx
	//log.Print("O", o)

	//log.Print("O", o)
	//	v.cibleX = 424
	//	v.cibleY = 240
	//log.Print("cibleXY")
	//log.Print(float64(v.cibleX)*v.cmPerPx, float64(v.cibleY)*v.cmPerPx)
	m := Point{x: (float64(v.cibleX) - o.x) * v.cmPerPx, y: (float64(v.cibleY) - o.y) * v.cmPerPx, z: dist / (float64(cible.W) / float64(visageA30cm.W))}
	o.x = o.x * v.cmPerPx
	o.y = o.y * v.cmPerPx
	o.z = o.z * v.cmPerPx
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
	pointM.z = rayon * sinZigm / v.cmPerPx
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
	pointM.x = rayon * cosTheta * cosZigm / v.cmPerPx

	ih := math.Sqrt(float64((h.x-i.x)*(h.x-i.x) + (h.y-i.y)*(h.y-i.y) + (h.z-i.z)*(h.z-i.z)))
	sinTheta := detY * ih / oh
	//log.Print("sinT", sinTheta)
	pointM.y = rayon * cosZigm * sinTheta / v.cmPerPx
	//log.Print("pM", pointM)
	return pointM
}
func (v *Visage) calculCible(cible *Cible) {
	x := cible.X
	y := cible.Y
	w := cible.W
	h := cible.H
	v.cibleX = (x + w/2) * v.size.W / v.capWidth
	v.cibleY = (y - h/2) * v.size.H / v.capHeight

	//visageA10cm := visage.Cible{X: 1, Y: 1, W: 100, H: 100}
	//cible := visage.Cible{X: 1, Y: 1, W: 30, H: 30}
	//log.Print("left")
	pointL := v.calcul(cible, Point{x: v.leftEllipseX, y: v.leftEllipseY, z: 0})
	//log.Print("right")

	pointR := v.calcul(cible, Point{x: v.rightEllipseX, y: v.rightEllipseY, z: 0})
	v.cible = cible
	v.cible.X = v.cible.X * v.size.W / v.capWidth
	v.cible.Y = v.cible.Y * v.size.W / v.capWidth
	v.cible.W = v.cible.W * v.size.W / v.capWidth
	v.cible.H = v.cible.H * v.size.W / v.capWidth

	//log.Println(pointM)
	v.leftX = pointL.x + v.leftEllipseX   //+ -(float64(v.size.W/2) - )
	v.leftY = pointL.y + v.leftEllipseY   //+ v.leftEllipseY
	v.rightX = pointR.x + v.rightEllipseX //+ (float64(v.size.W/2) - v.rightEllipseX)
	v.rightY = pointR.y + v.rightEllipseY //v.rightEllipseY +

}

/*
func (v *Visage) calculCible(cible *Cible) {

	x := cible.X
	y := cible.Y
	w := cible.W
	h := cible.H
	v.cibleX = (x + w/2) * v.size.W / v.capWidth
	v.cibleY = (y - h/2) * v.size.H / v.capHeight

	v.cible = cible
	v.cible.X = v.cible.X * v.size.W / v.capWidth
	v.cible.Y = v.cible.Y * v.size.W / v.capWidth
	v.cible.W = v.cible.W * v.size.W / v.capWidth
	v.cible.H = v.cible.H * v.size.W / v.capWidth

	//log.Print(v.leftEllipseX, v.leftEllipseY)
	m := (v.leftEllipseY - float64(v.cibleY)) / (v.leftEllipseX - float64(v.cibleX))
	//m := float64(v.cibleY) / float64(v.cibleX)

	b := visage.eyeWidth  // * 2 / 3
	a := visage.eyeHeight //* 2 / 3
	//log.Print(m, b, a)
	//X = (+ ou -) a*b / sqrt( b^2 + a^2*m^2
	v.X1 = a * b / math.Sqrt(b*b+a*a*m*m)
	v.X2 = -a * b / math.Sqrt(b*b+a*a*m*m)
	if float64(v.cibleX) < v.leftEllipseX {
		v.leftX = -a * b / math.Sqrt(b*b+a*a*m*m)
	} else {
		v.leftX = a * b / math.Sqrt(b*b+a*a*m*m)
	}
	v.leftY = m * v.leftX

	v.Y1 = m * v.X1
	v.Y2 = m * v.X2

	v.X1 += v.leftEllipseX
	v.X2 += v.leftEllipseX
	v.Y1 += v.leftEllipseY
	v.Y2 += v.leftEllipseY

	v.leftY += v.leftEllipseY
	v.leftX += v.leftEllipseX

	m = (v.rightEllipseY - float64(v.cibleY)) / (v.rightEllipseX - float64(v.cibleX))
	//m := float64(v.cibleY) / float64(v.cibleX)

	if float64(v.cibleX) < v.rightEllipseX {
		v.rightX = -a * b / math.Sqrt(b*b+a*a*m*m)
	} else {
		v.rightX = a * b / math.Sqrt(b*b+a*a*m*m)
	}
	v.rightY = m * v.rightX
	v.rightY += v.rightEllipseY
	v.rightX += v.rightEllipseX

	//log.Print(v.X1, v.Y1, v.X2, v.Y2)
}
*/
func (v *Visage) moveEyeLeft() {

	/*
		v.leftVy = Max(-1*v.maxVt, Min(v.maxVt,
			(v.leftTy-int(v.leftY))/v.ratioVt))
		v.leftVx = Max(-1*v.maxVt, Min(v.maxVt,
			(v.leftTx-int(v.leftX))/v.ratioVt))
		v.leftY += float64(v.leftVy)
		v.leftX += float64(v.leftVx)
		if v.leftX < float64(v.rad) {
			v.leftX = float64(v.rad)
		} else if v.leftX > float64(v.size.w-v.rad) {
			v.leftX = float64(v.size.w - v.rad)
		}
		if v.leftY < float64(v.rad) {
			v.leftY = float64(v.rad)
		} else if v.leftY > float64(v.size.h-v.rad) {
			v.leftY = float64(v.size.h - v.rad)
		}
	*/
}

func (v *Visage) moveEyeRight() { /*
		v.rightVy = Max(-1*v.maxVt, Min(v.maxVt,
			(v.rightTy-int(v.rightY))/v.ratioVt))
		v.rightVx = Max(-1*v.maxVt, Min(v.maxVt,
			(v.rightTx-int(v.rightX))/v.ratioVt))
		v.rightY += float64(v.rightVy)
		v.rightX += float64(v.rightVx)
		if v.rightX < float64(v.rad) {
			v.rightX = float64(v.rad)
		} else if v.rightX > float64(v.size.w-v.rad) {
			v.rightX = float64(v.size.w - v.rad)
		}
		if v.rightY < float64(v.rad) {
			v.rightY = float64(v.rad)
		} else if v.rightY > float64(v.size.h-v.rad) {
			v.rightY = float64(v.size.h - v.rad)
		}
	*/
}

func (v *Visage) composite(win *pixelgl.Window) {
	win.Clear(pixel.RGB(1, 1, 1))
	imd := imdraw.New(nil)
	//	imd.Color = pixel.RGB(1, 1, 1)

	//	imd.Push(pixel.V(float64(visage.leftEllipseX), visage.leftEllipseY),
	//		pixel.V(float64(visage.rightEllipseX), visage.rightEllipseY))

	//	imd.Ellipse(pixel.V(visage.eyeWidth, visage.eyeHeight), 0)

	//pupilles
	mat := pixel.IM
	mat = mat.Scaled(pixel.V(0, 0), float64(visage.eyeRadius)/visage.pd.Bounds().W())
	mat = mat.Moved(pixel.V(visage.leftX, visage.leftY))
	visage.pupil.Draw(win, mat)

	mat = pixel.IM
	mat = mat.Scaled(pixel.V(0, 0), float64(visage.eyeRadius)/visage.pd.Bounds().W())
	mat = mat.Moved(pixel.V(visage.rightX, visage.rightY))
	visage.pupil.Draw(win, mat)

	if visage.blink {
		//	log.Print("blink")
		imd.Color = pixel.RGB(0.996, 0.764, 0.674)
		imd.Push(pixel.V(float64(visage.leftEllipseX), visage.leftEllipseY+float64(visage.xxx)),
			pixel.V(float64(visage.rightEllipseX), visage.rightEllipseY+float64(visage.xxx)))
		imd.Ellipse(pixel.V(visage.eyeWidth, visage.eyeHeight), 0)

		//logging.info(visage.xxx)
		//	log.Print(visage.xxx)
		visage.xxx = visage.xxx + visage.xxS*5
		//		log.Print(visage.xxx)
		if visage.xxx > visage.maxBlink {
			visage.xxx = visage.maxBlink
			visage.xxS = -1 * visage.xxS
		}
		if visage.xxx < 0 {
			visage.xxx = 0
			visage.xxS = -1 * visage.xxS
		}
		if visage.xxx == visage.maxBlink && visage.xxS == -1 && visage.blink {
			visage.blink = false
		}
	}
	imd.Draw(win)

	//mask
	canvas := pixelgl.NewCanvas(win.Bounds())
	imd2 := imdraw.New(nil)
	canvas.SetBounds(win.Bounds())
	canvas.Clear(pixel.Alpha(0))
	canvas.SetComposeMethod(pixel.ComposeXor)

	imd2.Color = pixel.RGB(0.996, 0.764, 0.674)
	imd2.Push(pixel.V(0, 0), pixel.V(win.Bounds().W(), win.Bounds().H()))
	imd2.Rectangle(0)

	imd2.Color = pixel.RGB(0, 0, 0)
	imd2.Push(pixel.V(float64(visage.leftEllipseX), visage.leftEllipseY),
		pixel.V(float64(visage.rightEllipseX), visage.rightEllipseY))
	imd2.Ellipse(pixel.V(visage.eyeWidth, visage.eyeHeight), 0)

	imd2.Draw(canvas)
	canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	//point de cible
	imd3 := imdraw.New(nil)
	imd3.Color = pixel.RGB(0, 0, 1)
	imd3.Push(pixel.V(float64(visage.cibleX), float64(visage.cibleY)))
	imd3.Circle(5, 0)
	if visage.cible != nil {

		imd3.Push(pixel.V(float64(visage.cible.X), float64(visage.cible.Y)), pixel.V(float64(visage.cible.X+visage.cible.W), float64(visage.cible.Y-visage.cible.H)))
		imd3.Rectangle(2)
	}
	/*	imd.Color = pixel.RGB(1, 0, 0)
		imd.Push(pixel.V(float64(visage.X1), float64(visage.Y1)))
		imd.Circle(5, 0)
		imd.Color = pixel.RGB(0, 1, 0)
		imd.Push(pixel.V(float64(visage.X2), float64(visage.Y2)))
		imd.Circle(5, 0)
	*/
	//sourcils
	imd.Color = pixel.RGB(0, 0, 0)
	imd.Push(pixel.V(visage.leftEllipseX, visage.leftEllipseY+visage.eyeHeight), pixel.V(visage.rightEllipseX, visage.rightEllipseY+visage.eyeHeight))
	imd.EllipseArc(pixel.V(visage.eyeWidth, visage.eyeHeight), math.Pi/5, 4*math.Pi/5, 5)

	imd.Push(pixel.V(visage.leftEllipseX, visage.leftEllipseY), pixel.V(visage.rightEllipseX, visage.rightEllipseY))
	imd.EllipseArc(pixel.V(visage.eyeWidth, visage.eyeHeight), 0, 2*math.Pi, 1)

	imd.Draw(win)
	//imd3.Draw(win)
}
