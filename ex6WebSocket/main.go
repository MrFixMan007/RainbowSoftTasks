package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	channel := make(chan int)
	wg := &sync.WaitGroup{}
	fmt.Println(r.URL.Query().Get("flow"))
	fmt.Println(r.URL.Query().Get("limit"))
	fmt.Println(r.URL.Query().Get("count"))

	flow, _ := strconv.Atoi(r.URL.Query().Get("flow"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))

	if flow <= 0 || limit <= 0 || count <= 0 {
		connection.WriteMessage(websocket.TextMessage, []byte("Входные данные должны быть > 0"))
	} else {
		wg.Add(1)
		go Printer(flow, limit, count, wg, channel)
		for v := range channel {
			connection.WriteMessage(websocket.TextMessage, []byte{byte(v)})
			fmt.Println("отправил", v)
		}
		wg.Wait()
	}

	// for {
	// 	_, message, _ := connection.ReadMessage()

	// 	connection.WriteMessage(websocket.TextMessage, message)
	// 	go PrintMessage(string(message))
	// }
}

// Printer создаёт горутины, передаёт им параметры. Передаст в канал столько
// чисел, сколько требуется.
func Printer(flow int, limit int, count int, wg *sync.WaitGroup, out chan int) {
	defer fmt.Println("Printer завершил работу -")
	defer wg.Done()
	fmt.Println("Printer начал работу +")
	channel := make(chan int)
	quit := make(chan bool)
	wg1 := &sync.WaitGroup{}

	for i := 0; i < flow; i++ {
		wg1.Add(1)
		go Generator(limit, channel, quit, wg1, i)
	}

	for i := 0; i < int(count); i++ {
		quit <- true
		out <- <-channel
	}

	close(quit)
	wg1.Wait()
	close(channel)
	close(out)
}

// Generator создаёт случайные числа от 0 до лимита и передаёт в канал
func Generator(limit int, channel chan int, quit chan bool, wg1 *sync.WaitGroup, index int) {
	defer fmt.Printf("Generator %d завершился -\n", index+1)
	defer wg1.Done()
	fmt.Printf("Generator %d запустился +\n", index+1)

	for {
		_, ok := <-quit
		if ok {
			rand.Seed(time.Now().UnixNano())
			channel <- rand.Intn(int(limit))
		} else {
			break
		}
	}
}

func PrintMessage(message string) {
	fmt.Println(message)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var v = struct {
		Host string
	}{
		r.Host,
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
