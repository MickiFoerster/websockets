// echo-server is from
// https://godoc.org/golang.org/x/net/websocket#example-Handler
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

// echoServer serves websocket requests
func echoServer(ws *websocket.Conn) {
	defer func() {
		log.Println("connection handler exits")
	}()
	fmt.Println("local address:", ws.LocalAddr())
	fmt.Println("remote address:", ws.RemoteAddr())
	buf := make([]byte, 4096)
	n, _ := ws.Read(buf)
	log.Printf("Received %d bytes: %q\n", n, buf[:n])
	ws.Write([]byte("B"))
	n, _ = ws.Read(buf)
	log.Printf("Received %d bytes: %q\n", n, buf[:n])
	time.Sleep(2 * time.Second)
	ws.Write([]byte("C"))
	n, _ = ws.Read(buf)
	log.Printf("Received %d bytes: %q\n", n, buf[:n])
	time.Sleep(2 * time.Second)
	ws.Write([]byte("D"))
	n, _ = ws.Read(buf)
	log.Printf("Received %d bytes: %q\n", n, buf[:n])
	time.Sleep(2 * time.Second)
	ws.Close()
}

func handleMainRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle /")
	err := tpl.ExecuteTemplate(w, "tpl.gohtml", `This is a text`)
	if err != nil {
		log.Fatal(err)
	}
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
