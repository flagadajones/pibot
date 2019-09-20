package main

import (
	//"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

var maxBlink int = 40
var visage *Visage

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Eye",
		Bounds: pixel.R(0, 0, float64(visage.size.w), float64(visage.size.h)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	//	win.SetMonitor(pixelgl.PrimaryMonitor())
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Black)
	f, err := os.Open("./iris.png")
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
	log.Print("radiu")
	log.Print(visage.eyeRadius)
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
	x, y, w, h int
}

type CapSize struct {
	w, h int
}

func (v *Visage) draw(cible *Cible) {
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
}

func (v *Visage) Init(capSize *CapSize) {
	visage = v
	v.capWidth = capSize.w
	v.capHeight = capSize.h
	//v.size = CapSize{w: 848, h: 480}
	v.size = CapSize{w: 212, h: 120}

	v.maxVt = 50
	v.ratioVt = 2
	v.blink = true
	//	v.lastrect = nil

	v.starty = int(v.size.h / 2)
	v.startx = int(v.size.w / 2)
	v.cibleX = v.startx
	v.cibleY = v.starty

	v.eyeRadius = float64(v.size.h) / 4
	v.eyeRatio = 3.0 / 5.0
	v.eyeWidth = v.eyeRadius
	v.eyeHeight = float64(v.eyeRadius) * v.eyeRatio
	log.Print(v.starty)
	log.Print(v.eyeRadius)
	log.Print(v.eyeRatio)

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

	v.xxx = maxBlink
	v.xxS = 1
	//	v.blink = false

}
func (v *Visage) run() {
	pixelgl.Run(run)

}
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
	Xa := start.w - eye.w
	Ya := start.h - eye.h
	Xb := cible.w - eye.w
	Yb := cible.h - eye.h

	Na := math.Sqrt(float64(Xa*Xa + Ya*Ya))
	Nb := math.Sqrt(float64(Xb*Xb + Yb*Yb))
	C := float64(Xa*Xb+Ya*Yb) / (Na * Nb)
	S := (Xa*Yb - Ya*Xb)
	angle := sign(S) * math.Acos(C)
	return angle, Nb
}

func (v *Visage) calculCible(cible *Cible) {

	x := cible.x
	y := cible.y
	w := cible.w
	h := cible.h
	v.cibleX = x + w/2 + (v.size.w-v.capWidth)/2
	v.cibleY = y + h/2 + (v.size.h-v.capHeight)/2
	log.Print("calculCible", v.cibleX, v.cibleY)
	//cibleX := x + w/2 + (v.size.w-v.capWidth)/2
	// symetrie axiale
	//v.cibleX = int(2*v.startx - cibleX)
	//v.cibleY = int(y + h/2 + (v.size.h-v.capHeight)/2)
	// ratioDistance = w/212
	//ratioDistance := 1.0

	/*
		v.cibleX = int(2*v.startx - cibleX)

		angleLeft, left := v.calculRayon(
			CapSize{v.startx, v.starty}, CapSize{v.cibleX, v.cibleY}, CapSize{int(v.leftEyeX), int(v.leftEyeY)})

		angleRight, right := v.calculRayon(
			CapSize{v.startx, v.starty}, CapSize{v.cibleX, v.cibleY}, CapSize{int(v.rightEyeX), int(v.rightEyeY)})
		angleRight = angleRight + math.Pi

		log.Print("aa")
		log.Print(angleLeft)
		log.Print(angleRight)

		v.leftTx = int(v.leftEyeX + float64(v.eyeMoveRadiusWidth/2*int(
			math.Min(1, ratioDistance*left/(left+right))*math.Cos(angleLeft))))
		v.leftTy = int(v.leftEyeY + float64(v.eyeMoveRadius/2*int(
			math.Min(1, ratioDistance*left/(left+right))*math.Sin(angleLeft))))

		v.rightTx = int(v.rightEyeX + float64(v.eyeMoveRadiusWidth/2*int(
			math.Min(1, ratioDistance*right/(left+right))*math.Cos(angleRight))))
		v.rightTy = int(v.rightEyeY + float64(v.eyeMoveRadius/2*int(
			math.Min(1, ratioDistance*right/(left+right))*math.Sin(angleRight))))
	*/

	log.Print(v.leftEllipseX, v.leftEllipseY)
	m := (v.leftEllipseY - float64(v.cibleY)) / (v.leftEllipseX - float64(v.cibleX))
	//m := float64(v.cibleY) / float64(v.cibleX)

	b := visage.eyeWidth  // * 2 / 3
	a := visage.eyeHeight //* 2 / 3
	log.Print(m, b, a)
	//X = (+ ou -) a*b / sqrt( b^2 + a^2*m^2
	v.X1 = a * b / math.Sqrt(b*b+a*a*m*m)
	v.X2 = -a * b / math.Sqrt(b*b+a*a*m*m)
	if float64(v.cibleX) < v.leftEllipseX {
		v.leftX = -a * b / math.Sqrt(b*b+a*a*m*m)
	} else {
		v.leftX = a * b / math.Sqrt(b*b+a*a*m*m)
	}
	v.leftY = m * v.leftX
	//	if v.cibleX < int(v.leftEllipseX) {
	//		v.X1 = -v.X1
	//	}
	v.Y1 = m * v.X1
	v.Y2 = m * v.X2
	//	if v.cibleY < int(v.leftEllipseY) {
	//		v.Y1 = -v.Y1
	//	}
	//	v.leftY = v.Y1 + v.leftEllipseY
	//	v.leftX = v.X1 + v.leftEllipseX

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

	log.Print(v.X1, v.Y1, v.X2, v.Y2)
}

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

	/*


	   // #   point = pygame.draw.circle(self.display, (255, 0, 0),(int(self.leftTx),int(self.leftTy)) , 10)
	   // #   point1 = pygame.draw.circle(self.display, (255, 0, 0),(int(self.rightTx),int(self.rightTy)) , 10)

	   // #   pygame.draw.circle(self.display, (255, 255, 0),(int(self.leftEyeX),int(self.leftEyeY)) , 10)

	   	if(self.blink):
	   		p2 = pygame.draw.ellipse(self.display, (254, 195, 172), [
	   			self.leftEllipseX, self.leftEllipseY-self.xxx, self.eyeWidth, self.eyeHeight], 0)
	   		p3 = pygame.draw.ellipse(self.display, (254, 195, 172), [
	   			self.rightEllipseX, self.rightEllipseY-self.xxx, self.eyeWidth, self.eyeHeight], 0)
	   		logging.info(self.xxx)
	   		self.xxx = self.xxx+self.xxS*100
	   		if(self.xxx > 150):
	   			self.xxx = 150
	   			self.xxS = -1*self.xxS
	   		if(self.xxx < 0):
	   			self.xxx = 0
	   			self.xxS = -1*self.xxS
	   		if(self.xxx == 150 and self.xxS == -1 and self.blink):
	   			self.blink = False


	*/
	//  	point2 = pygame.draw.circle(self.display, (0, 0, 255), (int(self.cibleX), int(self.cibleY)), 10)

	/*

			if(self.blink):
				pygame.display.update(
					[p2, p3])
		//#  pygame.display.update([point,point1,point2,point3,point4,leftPupil,leftEye,rightPupil,rightEye,self.lastrect])
		//   # pygame.display.update()
	*/

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
		log.Print("blink")
		imd.Color = pixel.RGB(0.996, 0.764, 0.674)
		imd.Push(pixel.V(float64(visage.leftEllipseX), visage.leftEllipseY+float64(visage.xxx)),
			pixel.V(float64(visage.rightEllipseX), visage.rightEllipseY+float64(visage.xxx)))
		imd.Ellipse(pixel.V(visage.eyeWidth, visage.eyeHeight), 0)

		//logging.info(visage.xxx)
		log.Print(visage.xxx)
		visage.xxx = visage.xxx + visage.xxS*5
		log.Print(visage.xxx)
		if visage.xxx > maxBlink {
			visage.xxx = maxBlink
			visage.xxS = -1 * visage.xxS
		}
		if visage.xxx < 0 {
			visage.xxx = 0
			visage.xxS = -1 * visage.xxS
		}
		if visage.xxx == maxBlink && visage.xxS == -1 && visage.blink {
			//visage.blink = false
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
	imd.Color = pixel.RGB(0, 0, 1)
	imd.Push(pixel.V(float64(visage.cibleX), float64(visage.cibleY)))
	imd.Circle(5, 0)
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

}
