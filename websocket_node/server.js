var fs = require("fs");
var http = require("http");
var WebSocketServer = require("websocket").server;

var httpServer = http.createServer(function(req,res) {
  fs.readFile(__dirname+"/index.html", function(error, data) {
    if(error) {
      res.writeHead(500);
      return res.end('Fehler beim Laden der Datei index.html');
    }
    else {
      res.writeHead(200);
      res.end(data);
    }
  })
});

var webSocketServer = new WebSocketServer({
  httpServer:httpServer,
  autoAcceptConnections:true,
  maxReceivedFrameSize: 64*1024*1024,
  maxReceivedMessageSize: 64*1024*1024,
});

webSocketServer.on("connect",function(socket){
  socket.on("message",function (message){
    if(message.type=="binary"){
      socket.sendBytes(message.binaryData);
    }
  });
});

httpServer.listen(12345);
console.log("Der EchoServer läuft auf dem Port ", httpServer.address().port);
