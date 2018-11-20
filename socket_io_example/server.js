var fs = require('fs'),
  http = require('http');

var httpServer = http.createServer(function(request, response){
  fs.readFile(__dirname+"/index.html", function(error, data){
    if(error) {
      response.writeHead(500);
      return response.end('Fehler beim Laden der Datei index.html');
    }
    else {
      response.writeHead(200);
      response.end(data);
    }
  });
});

var io = require('socket.io').listen(httpServer);

io.sockets.on('connection', function(socket){
  socket.on('message', function(message){
    socket.send(message);
  })

  socket.on('myevent', function(message){
    socket.emit('myevent', "Eine Nachricht von myevent: " + message);
  })
})

httpServer.listen(12345);
console.log("Der EchoServer läuft auf dem Port", httpServer.address().port);
