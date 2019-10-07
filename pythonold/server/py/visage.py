import pygame
from pygame.locals import *
import numpy as np
import logging
import random
from math import sqrt, asin, cos, sin, acos, pi
import os
dir_path = os.path.dirname(os.path.realpath(__file__))
pygame.init()


class Visage(object):

    def draw(self, cible):
        self.calcul_cible(cible)
        self.move_eyeLeft()
        self.move_eyeRight()
        self.composite()

    def __init__(self, capSize):
        self.capWidth = capSize[0]
        self.capHeight = capSize[1]
#        self.size = (320, 200)
#         self.size = (1280, 720)
        self.size = (848, 480)
        self.display = pygame.display.set_mode(self.size)
#        self.display = pygame.display.set_mode(self.size, pygame.FULLSCREEN)
        logging.warn(pygame.display.list_modes())
        self.display.fill((0, 0, 0))
        pygame.display.set_caption('Eye')

        self.maxVt = 50
        self.ratioVt = 2

        self.lastrect = None
        self.starty = int(self.size[1]/2)
        self.startx = int(self.size[0]/2)

        self.eyeRadius = int(self.size[1]/4)
        self.eyeRatio = 3/5
        self.eyeWidth = int(self.eyeRadius*2)
        self.eyeHeight = int(self.eyeRadius*self.eyeRatio*2)

        self.pupil = pygame.image.load(dir_path+'/iris.png')
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
        self.blink = False

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

    def calcul_cible(self, cible):

        x = cible["x"]
        y = cible["y"]
        w = cible["w"]
        h = cible["h"]

        cibleX = x+w/2+(self.size[0]-self.capWidth)/2
        # symetrie axiale
        self.cibleX = int(2*self.startx-cibleX)
        self.cibleY = int(y+h/2+(self.size[1]-self.capHeight)/2)
        # ratioDistance = w/212
        ratioDistance = 1

        self.cibleX = int(2*self.startx-cibleX)

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

        mask.fill((255, 255, 255, 255))

        pygame.draw.ellipse(mask, (0, 0, 0, 0), [
                            self.leftEllipseX, self.leftEllipseY, self.eyeWidth, self.eyeHeight], 0)

        pygame.draw.ellipse(mask, (0, 0, 0, 0), [
                            self.rightEllipseX, self.rightEllipseY, self.eyeWidth, self.eyeHeight], 0)

        full = pygame.Surface(self.size, pygame.SRCALPHA)

        full.fill((254, 195, 172))

        full.blit(mask, (0, 0), None, pygame.BLEND_RGBA_MULT)

        self.display.blit(full, (0, 0), None)

        point2 = pygame.draw.circle(
            self.display, (0, 0, 255), (int(self.cibleX), int(self.cibleY)), 10)
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
