package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	start()
}

// homePage открывает шаблон index.html
func homePage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("index.html")
	templ.Execute(w, nil)
}

// out открывает шаблон index.html
func out(w http.ResponseWriter, r *http.Request) {
	var answer string
	channel := make(chan int)
	wg := &sync.WaitGroup{}

	flow, _ := strconv.Atoi(r.FormValue("flow"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	count, _ := strconv.Atoi(r.FormValue("count"))

	if flow <= 0 || limit <= 0 || count <= 0 {
		answer = "Входные данные должны быть > 0"
	} else {
		wg.Add(1)
		go Printer(flow, limit, count, wg, channel)
		for v := range channel {
			answer += fmt.Sprintf("%d ", v)
		}
		wg.Wait()
	}

	templ, _ := template.ParseFiles("out.html")
	templ.Execute(w, answer)
}

// start запускает сервер, открывает порт
func start() {
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles/"))))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/out", out)
	http.ListenAndServe(":4434", nil)
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
