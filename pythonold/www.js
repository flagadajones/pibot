#!/usr/bin/env node

/**
 * Module dependencies.
 */

var app = require('./server/app');
var debug = require('debug')('my-app:server');
var http = require('http');

/**
 * Get port from environment and store in Express.
 */

var port = normalizePort(process.env.PORT || '3000');
app.set('port', port);

/**
 * Create HTTP server.
 */

var server = http.createServer(app);

/**
 * Socket.io
 */

app.io.attach(server);


/**
 * Listen on provided port, on all network interfaces.
 */

server.listen(port);
server.on('error', onError);
server.on('listening', onListening);

/**
 * Start Python 
 * 
 */

var spawn = require('child_process').spawn,
  py = spawn('python', [String(__dirname) + './server/py/main.py']);
py.stderr.on('data', (data) => {
  console.log(`main.js -> stderr: ${data}`);
});

function exitHandler(data, signal) {
  console.log('main.js -> success code:' + data + ' ' + signal);
}
function errorHandler(data) {
  console.log('train.js -> error code:' + data);
}
py.addListener('close', exitHandler);
py.addListener('error', errorHandler);

/**
 * Normalize a port into a number, string, or false.
 */

function normalizePort(val) {
  var port = parseInt(val, 10);

  if (isNaN(port)) {
    // named pipe
    return val;
  }

  if (port >= 0) {
    // port number
    return port;
  }

  return false;
}

/**
 * Event listener for HTTP server "error" event.
 */

function onError(error) {
  if (error.syscall !== 'listen') {
    throw error;
  }

  var bind = typeof port === 'string'
    ? 'Pipe ' + port
    : 'Port ' + port;

  // handle specific listen errors with friendly messages
  switch (error.code) {
    case 'EACCES':
      console.error(bind + ' requires elevated privileges');
      process.exit(1);
      break;
    case 'EADDRINUSE':
      console.error(bind + ' is already in use');
      process.exit(1);
      break;
    default:
      throw error;
  }
}

/**
 * Event listener for HTTP server "listening" event.
 */

function onListening() {
  var addr = server.address();
  var bind = typeof addr === 'string'
    ? 'pipe ' + addr
    : 'port ' + addr.port;
  debug('Listening on ' + bind);
  // Log on console to confirm correct execution
  console.log(`Server listening on port: ${port}.`);
}