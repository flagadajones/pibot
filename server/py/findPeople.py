import cv2
from ocv_detection import *
from ocv_recognition import *
import os
dir_path = os.path.dirname(os.path.realpath(__file__))


class FindPeople(object):

    def reload_train(self):
        self.recognizer.load_trainset()
        self.recognizer.train()

    def __init__(self, capSize):
        self.capSize = capSize
        self.capture = cv2.VideoCapture(0)
        self.capture.set(cv2.CAP_PROP_FRAME_WIDTH, self.capSize[0])
        self.capture.set(cv2.CAP_PROP_FRAME_HEIGHT, self.capSize[1])
        self.detector = OpenCVFaceFrontalDetection(
            archive_folder=dir_path+"/../../public/archives/",
            debug=False)
        self.recognizer = OpenCVFaceRecognitionLBPH(dir_path+"/../../public/faces/",
                                                    archive_folder=dir_path+"/../../public/archives/")
        self.recognizer.load_trainset()
        self.recognizer.train()

    def archive_items_frames(self):
        logging.debug('archive_items_frames')
        self.detector.archive_items_frames()

    def extract_items_frames(self):
        logging.debug('extract_items_frames')
        self.detector.extract_items_frames()

    def find_faces(self):
        ret, img = self.capture.read()
        # self.gray = cv2.cvtColor(img,  cv2.COLOR_BGR2GRAY)
        #cv2.imshow('img', img)
        self.detector.set_frame(img)
        self.detector.find_items()

     #   T.extract_items_frames()
     #   T.archive_items_frames()
     #   T.archive_with_items()
     #   logging.info(T.items)

        self.detector.extract_items_frames()
        result = []
        for item in self.detector.get_items_frames(grayscale=False):
            known, identity, confidence = self.recognizer.recognize(
                item["frame"])
            label = "{0} ({1})".format(identity, confidence)
            logging.warn("Trouvé : {0}".format(label))
            x = item["x"]
            y = item["y"]
            self.detector.add_label(label, x, y)
            if known:
                result.append(
                    {'x': int(item["x"]), 'y': int(item["y"]), 'w': int(item["w"]), 'h': int(item["h"]), 'label': identity})
            if not known:
                pass
                # TODO : save in a unknown folder
                # TODO : save in a unknown folder
                # TODO : save in a unknown folder
         # self.detector.archive_items_frames()
          #  self.detector.archive_with_items()

        return result
    # for (x, y, w, h) in T.items:
