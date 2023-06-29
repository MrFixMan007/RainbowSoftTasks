package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serveWs запустился")
	defer fmt.Println("serveWs закрылся")
	fmt.Println(w)
	fmt.Println(r)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		fmt.Println(err)
		return
	}
	ws.WriteMessage(websocket.PingMessage, []byte{13, 21, 32, 13, 213})
	fmt.Println("serveWs написал")
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		"Hi!!!!!",
	}
	templ, _ := template.ParseFiles("index.html")
	templ.Execute(w, &v)
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
