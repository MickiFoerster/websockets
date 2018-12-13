// echo-server is from
// https://godoc.org/golang.org/x/net/websocket#example-Handler
package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	fmt.Println("local address:", ws.LocalAddr())
	fmt.Println("remote address:", ws.RemoteAddr())
	io.Copy(ws, ws)
}

// This example demonstrates a trivial echo server.
func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
