import pygame
import numpy as np
from ocv_detection import *
from pygame.locals import *
import logging
import random
import cv2
from math import sqrt, asin, cos, sin, acos, pi
pygame.init()


class EvilEye(object):
    def __init__(self):
        self.size = (848, 480)
        self.display = pygame.display.set_mode(self.size)
#        self.display = pygame.display.set_mode(self.size, pygame.FULLSCREEN)
        self.display.fill((0, 0, 0))
        pygame.display.set_caption('Eye')

        self.clock = pygame.time.Clock()
        self.cameratick = pygame.time.get_ticks()

        self.cascade = cv2.CascadeClassifier(
            'haarcascade_frontalface_default.xml')
        self.capture = cv2.VideoCapture(0)
        self.capWidth = 1024
        self.capHeight = 768
        self.capture.set(cv2.CAP_PROP_FRAME_WIDTH, self.capWidth)
        self.capture.set(cv2.CAP_PROP_FRAME_HEIGHT, self.capHeight)

#    self.gray = cv2.cvCreateImage(cv2.cvSize(640,480), 8, 1)
        self.maxVt = 50
        self.ratioVt = 2

        self.lastrect = None
        self.starty = int(self.size[1]/2)
        self.startx = int(self.size[0]/2)

        self.eyeRadius = int(self.size[1]/4)
        self.eyeRatio = 3/5
        self.eyeWidth = int(self.eyeRadius*2)
        self.eyeHeight = int(self.eyeRadius*self.eyeRatio*2)

        self.pupil = pygame.image.load('iris.png')
        self.pupil = self.pupil.convert_alpha()
        self.pupil = pygame.transform.scale(
            self.pupil, [self.eyeRadius, self.eyeRadius])
        self.rad = self.pupil.get_width()/2

        self.eyeMoveRadius = self.eyeHeight * 2/3
        self.eyeMoveRadiusWidth = self.eyeWidth * 2/3
        self.cibleX = self.startx
        self.cibleY = self.starty

        self.leftEllipseX = int(self.startx - self.startx/2-self.eyeRadius)
        self.leftEllipseY = int(self.starty-self.eyeRadius*self.eyeRatio)
        self.leftEyeX = int(self.leftEllipseX+self.eyeWidth/2)
        self.leftEyeY = int(self.leftEllipseY+self.eyeHeight/2)
        self.leftX = self.leftEyeX
        self.leftY = self.leftEyeY
        self.leftTx = self.leftX
        self.leftTy = self.leftY
        self.leftVx = 0
        self.leftVy = 0

        self.rightEllipseX = int(self.startx + self.startx/2-self.eyeRadius)
        self.rightEllipseY = int(self.starty-self.eyeRadius*self.eyeRatio)
        self.rightEyeX = int(self.rightEllipseX+self.eyeWidth/2)
        self.rightEyeY = int(self.rightEllipseY+self.eyeHeight/2)
        self.rightX = self.rightEyeX
        self.rightY = self.rightEyeY
        self.rightTx = self.rightX
        self.rightTy = self.rightY
        self.rightVx = 0
        self.rightVy = 0

        self.xxx = 150
        self.xxS = 1
        self.blink = True

    def calcul_rayon(self, start, cible, eye):
        Xa = start[0] - eye[0]
        Ya = start[1] - eye[1]
        Xb = cible[0] - eye[0]
        Yb = cible[1] - eye[1]

        Na = sqrt(Xa*Xa+Ya*Ya)
        Nb = sqrt(Xb*Xb+Yb*Yb)
        C = (Xa*Xb+Ya*Yb)/(Na*Nb)
        S = (Xa*Yb-Ya*Xb)
        angle = np.sign(S)*acos(C)
        return (angle, Nb)

    def find_face(self):
        ret, img = self.capture.read()
        self.gray = cv2.cvtColor(img,  cv2.COLOR_BGR2GRAY)
        T = OpenCVFaceFrontalDetection(self.gray,
                                       archive_folder="../public/archives/",
                                       debug=False)
        T.find_items()
     #   T.extract_items_frames()
     #   T.archive_items_frames()
     #   T.archive_with_items()
     #   logging.info(T.items)
        '''
        self.gray = cv2.cvtColor(img,  cv2.COLOR_BGR2GRAY)
        faces = self.cascade.detectMultiScale(self.gray, 1.3, 5)
        '''
        for (x, y, w, h) in T.items:
            cv2.rectangle(img, (x, y), (x+w, y+h), (255, 0, 0), 2)
            cibleX = x+w/2+(self.size[0]-self.capWidth)/2
            # symetrie axiale
            self.cibleX = int(2*self.startx-cibleX)
            self.cibleY = int(y+h/2+(self.size[1]-self.capHeight)/2)

            # ratioDistance = w/212
            ratioDistance = 1

            angleLeft, left = self.calcul_rayon(
                (self.startx, self.starty), (self.cibleX, self.cibleY), (self.leftEyeX, self.leftEyeY))

            angleRight, right = self.calcul_rayon(
                (self.startx, self.starty), (self.cibleX, self.cibleY), (self.rightEyeX, self.rightEyeY))
            angleRight = angleRight+pi

            self.leftTx = self.leftEyeX + self.eyeMoveRadiusWidth/2 * \
                min(1, ratioDistance*left/(left+right)) * cos(angleLeft)
            self.leftTy = self.leftEyeY + self.eyeMoveRadius/2 * \
                min(1, ratioDistance*left/(left+right)) * sin(angleLeft)

            self.rightTx = self.rightEyeX + self.eyeMoveRadiusWidth/2 * \
                min(1, ratioDistance*right/(left+right)) * cos(angleRight)
            self.rightTy = self.rightEyeY + self.eyeMoveRadius/2 * \
                min(1, ratioDistance*right/(left+right)) * sin(angleRight)
            break

        cv2.imshow('img', img)

    def move_eyeLeft(self):
        self.leftVy = max(-1*self.maxVt, min(self.maxVt,
                                             (self.leftTy-self.leftY)/self.ratioVt))
        self.leftVx = max(-1*self.maxVt, min(self.maxVt,
                                             (self.leftTx-self.leftX)/self.ratioVt))
        self.leftY += int(self.leftVy)
        self.leftX += int(self.leftVx)
        if self.leftX < self.rad:
            self.leftX = self.rad
        elif self.leftX > self.size[0]-self.rad:
            self.leftX = self.size[0]-self.rad
        if self.leftY < self.rad:
            self.leftY = self.rad
        elif self.leftY > self.size[1]-self.rad:
            self.leftY = self.size[1]-self.rad

    def move_eyeRight(self):
        self.rightVy = max(-1*self.maxVt, min(self.maxVt,
                                              (self.rightTy-self.rightY)/self.ratioVt))
        self.rightVx = max(-1*self.maxVt, min(self.maxVt,
                                              (self.rightTx-self.rightX)/self.ratioVt))
        self.rightY += int(self.rightVy)
        self.rightX += int(self.rightVx)
        if self.rightX < self.rad:
            self.rightX = self.rad
        elif self.rightX > self.size[0]-self.rad:
            self.rightX = self.size[0]-self.rad
        if self.rightY < self.rad:
            self.rightY = self.rad
        elif self.rightY > self.size[1]-self.rad:
            self.rightY = self.size[1]-self.rad

    def composite(self):
        self.lastrect = self.display.fill((254, 195, 172), self.lastrect)
        leftEye = pygame.draw.ellipse(self.display, (255, 255, 255), [
                                      self.leftEllipseX, self.leftEllipseY, self.eyeWidth, self.eyeHeight], 0)

        rightEye = pygame.draw.ellipse(self.display, (255, 255, 255), [
                                       self.rightEllipseX, self.rightEllipseY, self.eyeWidth, self.eyeHeight], 0)

#    point3 = pygame.draw.ellipse(self.display, (0, 255, 0), [int(self.startx -self.startx/2-self.eyeMoveRadiusWidth/2),int(self.starty -self.eyeMoveRadius/2), int(self.eyeMoveRadiusWidth) , int(self.eyeMoveRadius)],0)
#    point4 = pygame.draw.ellipse(self.display, (0, 255, 0), [int(self.startx +self.startx/2-self.eyeMoveRadiusWidth/2),int(self.starty -self.eyeMoveRadius/2), int(self.eyeMoveRadiusWidth) , int(self.eyeMoveRadius)],0)

        leftPupil = self.display.blit(
            self.pupil, (self.leftX-self.rad, self.leftY-self.rad))
        rightPupil = self.display.blit(
            self.pupil, (self.rightX-self.rad, self.rightY-self.rad))

     #   point = pygame.draw.circle(self.display, (255, 0, 0),(int(self.leftTx),int(self.leftTy)) , 10)
     #   point1 = pygame.draw.circle(self.display, (255, 0, 0),(int(self.rightTx),int(self.rightTy)) , 10)

     #   pygame.draw.circle(self.display, (255, 255, 0),(int(self.leftEyeX),int(self.leftEyeY)) , 10)

        point2 = pygame.draw.circle(
            self.display, (0, 0, 255), (int(self.cibleX), int(self.cibleY)), 10)
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
        mask = pygame.Surface(self.size, pygame.SRCALPHA)
        mask.fill((254, 195, 172))
        pygame.draw.ellipse(mask, (0, 0, 0), [
                            self.leftEllipseX, self.leftEllipseY, self.eyeWidth, self.eyeHeight], 0)
        pygame.draw.ellipse(mask, (0, 0, 0), [
                            self.rightEllipseX, self.rightEllipseY, self.eyeWidth, self.eyeHeight], 0)

        self.display.blit(mask, (0, 0), None, pygame.BLEND_RGB_MAX)

        pygame.draw.arc(self.display, (0, 0, 0), [
                        self.leftEllipseX, self.leftEllipseY-50, self.eyeWidth, self.eyeHeight], pi/5, 4*pi/5, 10)

        pygame.draw.arc(self.display, (0, 0, 0), [
                        self.leftEllipseX, self.leftEllipseY, self.eyeWidth, self.eyeHeight], 0, 2*pi, 1)

        pygame.draw.arc(self.display, (0, 0, 0), [
                        self.rightEllipseX, self.rightEllipseY-50, self.eyeWidth, self.eyeHeight], pi/5, 4*pi/5, 10)

        pygame.draw.arc(self.display, (0, 0, 0), [
                        self.rightEllipseX, self.rightEllipseY, self.eyeWidth, self.eyeHeight], 0, 2*pi, 1)

        pygame.display.update(
            [point2, leftPupil, leftEye, rightPupil, rightEye, self.lastrect])

        if(self.blink):
            pygame.display.update(
                [p2, p3])
    #  pygame.display.update([point,point1,point2,point3,point4,leftPupil,leftEye,rightPupil,rightEye,self.lastrect])
       # pygame.display.update()

    def main(self):
        going = True
        pygame.time.set_timer(USEREVENT+1, random.randint(4000, 8000))
        while going:
            events = pygame.event.get()
            for e in events:
                if e.type == (USEREVENT+1):
                    self.blink = True
                    pygame.time.set_timer(USEREVENT+1, 0)
                    pygame.time.set_timer(
                        USEREVENT+1, random.randint(4000, 8000))
                if e.type == QUIT or (e.type == KEYDOWN and e.key == K_ESCAPE):
                    going = False

            now = pygame.time.get_ticks()
            if now > self.cameratick + 100:
                self.find_face()
                self.cameratick = now

            self.move_eyeLeft()
            self.move_eyeRight()
            self.composite()
         #   self.clock.tick(30)

            self.clock.tick()


logging.basicConfig(level=logging.INFO)
EvilEye().main()
