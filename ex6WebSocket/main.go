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

// start запускает сервер, открывает порт
func main() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":4444", nil)
}

// homePage открывает шаблон index.html
func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host string
	}{
		r.Host,
	}
	templ, _ := template.ParseFiles("index.html")
	templ.Execute(w, &v)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(conn, err)
		return
	}
	defer conn.Close()
	fmt.Println("upgrader.Upgrade")

}
