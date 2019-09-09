import pygame
from pygame.locals import *
import logging
from visage import *
import threading
from findPeople import *
from socketIO_client import SocketIO
import logging
logging.getLogger('requests').setLevel(logging.WARNING)
logging.basicConfig(level=logging.WARNING)


BLINKEVENT = USEREVENT+1
TRAINEVENT = USEREVENT+2
RELOADTRAINEVENT = USEREVENT+3

# def on_aaa_response(*args):
#    logging.warn('on_aaa_response')


class Main(object):

    def on_blink(self, *args):
        ev = pygame.event.Event(BLINKEVENT, {})
        pygame.event.post(ev)
        logging.warn('on_blink_response')

    def on_train(self, *args):
        ev = pygame.event.Event(TRAINEVENT, {})
        pygame.event.post(ev)
        logging.warn('on_train_response')

    def on_reloadTrain(self, *args):
        ev = pygame.event.Event(RELOADTRAINEVENT, {})
        pygame.event.post(ev)
        logging.warn('on_reloadTrain_response')

    def wait(self):
        while True:
         #   logging.warn("wait")
            self.socketIO.wait(seconds=1)

    def __init__(self):
        logging.warn("eeee")
        self.clock = pygame.time.Clock()
        self.cameratick = pygame.time.get_ticks()
        self.trainCameratick = self.cameratick

        self.train = False

        self.capSize = (1024, 768)
        self.visage = Visage(self.capSize)
        self.findPeople = FindPeople(self.capSize)
        self.socketIO = SocketIO('localhost', 3000)
        self.socketIO.on('blink', self.on_blink)
        self.socketIO.on('train', self.on_train)
        self.socketIO.on('reloadTrain', self.on_reloadTrain)

    def main(self):
        try:
            t = threading.Thread(target=self.wait)
            t.daemon = True
            t.start()
        except Exception as e:
            logging.error("Error: unable to start thread")
            logging.error(str(e))

        cible = {'x': 0, 'y': 0, 'w': 0, 'h': 0}
        going = True
        pygame.time.set_timer(BLINKEVENT, random.randint(4000, 8000))

        while going:
           # self.wait()
            events = pygame.event.get()
            now = pygame.time.get_ticks()
            for e in events:
                if e.type == (BLINKEVENT):
                    self.visage.blink = True
                    pygame.time.set_timer(BLINKEVENT, 0)
                    pygame.time.set_timer(
                        BLINKEVENT, random.randint(4000, 8000))
                if e.type == (TRAINEVENT):
                    logging.warn("start train")
                    self.train = True
                    self.trainTick = now
                    self.trainCameratick = now
                if e.type == (RELOADTRAINEVENT):
                    logging.warn("start reload")
                    self.findPeople.reload_train()

                if e.type == QUIT or (e.type == KEYDOWN and e.key == K_ESCAPE):

                    going = False

            if self.train == True and now > self.trainTick + 10000:
                logging.warn("fin train")
                self.train = False
            if now > self.cameratick + 100:

               # self.socketIO.emit("aaa", "tt")

                faces = self.findPeople.find_faces()

                if self.train == True and now > self.trainCameratick + 300:
                    logging.warn("trainnn")
                    self.findPeople.extract_items_frames()
                    self.findPeople.archive_items_frames()

                    self.trainCameratick = now
                # println(faces)
                if(len(faces) > 0):
                    cible = faces[0]
                    self.socketIO.emit("faces", faces)

                self.cameratick = now

            self.visage.draw(cible)
        #   self.clock.tick(30)

            # self.socketIO.wait(seconds=0)
            # self.socketIO.wait(seconds=1)
            self.clock.tick()


import os
dir_path = os.path.dirname(os.path.realpath(__file__))
logging.info(dir_path)

Main().main()
