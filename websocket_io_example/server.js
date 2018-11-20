var ws = require("websocket.io"),
  fs = require('fs'),
  http = require('http');

var httpServer = http.createServer(function (request,response) {
  fs.readFile (__dirname+"/index.html", function(error, data) {
    if(error) {
      response.writeHead(500);
      return response.end('Fehler beim Laden der Datei index.html');
    }
    else {
      response.writeHead(200);
      response.end(data);
    }
  })
});

// Vereint Webserver und WebSocket-Server
var wsServer = ws.attach(httpServer);

wsServer.on('connection', function(client) {
  client.on('message', function(message) {
    client.send(message);
  });
})

httpServer.listen(12345);
console.log("Der Echo-Server laeuft auf dem Port",
  httpServer.address().port);
