package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {
	origin := "http://192.168.0.15/"
	url := "ws://192.168.0.15:12345/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
