from socketIO_client import SocketIO
import logging
logging.getLogger('requests').setLevel(logging.WARNING)
logging.basicConfig(level=logging.DEBUG)


def on_aaa_response(*args):
    print('on_aaa_response', args)


socketIO = SocketIO('localhost', 3000)
socketIO.on('aaa_response', on_aaa_response)
while True:
    socketIO.emit('aaa', "tot")
    socketIO.wait(seconds=1)
