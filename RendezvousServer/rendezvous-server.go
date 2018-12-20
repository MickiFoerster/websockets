// echo-server is from
// https://godoc.org/golang.org/x/net/websocket#example-Handler
package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

// echoServer serves websocket requests
func echoServer(ws *websocket.Conn) {
	defer func() {
		log.Println("ws connection handler exits")
	}()
	log.Println("local address:", ws.LocalAddr())
	log.Println("local address:", ws.LocalAddr().Network())
	log.Println("remote address:", ws.RemoteAddr())
	log.Println("remote address:", ws.RemoteAddr().Network())
	log.Println("IsClientConn:", ws.IsClientConn())
	log.Println("IsServerConn:", ws.IsServerConn())
	log.Println("Location:", ws.Config().Location)
	log.Println("Origin:", ws.Config().Origin)
	log.Println("Protocol:", ws.Config().Protocol)
	log.Println("Version:", ws.Config().Version)
	log.Println("TlsConfig:", ws.Config().TlsConfig)
	log.Println("Header:", ws.Config().Header)
	buf := make([]byte, 4096)
Loop:
	for {
		n, _ := ws.Read(buf)
		log.Printf("Received %d bytes: %q\n", n, buf[:n])
		req := string(buf[:n])
		req = strings.ToUpper(req)
		switch req {
		case "HELLO":
			ws.Write([]byte("Hello"))
			break
		default:
			if strings.HasPrefix(req, "FROM ") {
				name := req[len("FROM "):]
				log.Printf("Client %q connected\n", name)
				break
			}
			log.Println("Protocol error: ", req)
			ws.Close()
			break Loop
		}
	}
}

func handleMainRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("handle / (client address: ", r.RemoteAddr, ")")
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
