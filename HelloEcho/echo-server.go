// echo-server is from
// https://godoc.org/golang.org/x/net/websocket#example-Handler
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// echoServer serves websocket requests
func echoServer(ws *websocket.Conn) {
	fmt.Println("local address:", ws.LocalAddr())
	fmt.Println("remote address:", ws.RemoteAddr())
	buf := make([]byte, 4096)
	n, err := ws.Read(buf)
	if err != nil {
		fmt.Println("Could not read from websocket", err)
	}
	log.Printf("Received %d bytes: %q\n", n, buf[:n])
	ws.Write([]byte("A"))
}

func handleMainRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handled")
	fc, err := ioutil.ReadFile("echo-client.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(fc)
}

// This example demonstrates a trivial echo server.
func main() {
	http.HandleFunc("/", handleMainRoute)
	http.Handle("/echo", websocket.Handler(echoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
