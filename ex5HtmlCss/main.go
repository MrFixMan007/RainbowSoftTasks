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

// homePage открывает шаблон index.html
func homePage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("index.html")
	templ.Execute(w, nil)
	// fmt.Fprintf(w, StyleAdress)
}

// out открывает шаблон index.html
func out(w http.ResponseWriter, r *http.Request) {
	var answer string
	flow, _ := strconv.Atoi(r.FormValue("flow"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	count, _ := strconv.Atoi(r.FormValue("count"))
	if flow <= 0 || limit <= 0 || count <= 0 {
		answer = "Входные данные должны быть > 0"
	} else {
		//wg := &sync.WaitGroup{}
		//wg.Add(1)
		arr := Printer(flow, limit, count) //TODO
		for _, v := range arr {
			answer += fmt.Sprintf("%d ", v)
		}
		//wg.Wait()
	}

	templ, _ := template.ParseFiles("out.html")
	templ.Execute(w, answer)
	// fmt.Fprintf(w, StyleAdress)
}

// start запускает сервер, открывает порт
func start() {
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles/"))))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/out", out)
	http.ListenAndServe(":4434", nil)
}

func main() {
	start()
}

// Printer создаёт горутины, передаёт им параметры. Выведет столько чисел, количество которых было передано.
func Printer(flow int, limit int, count int) []int {
	channel := make(chan int)
	quit := make(chan bool)
	wg1 := &sync.WaitGroup{}
	arr := []int{}
	for i := 0; i < flow; i++ {
		wg1.Add(1)
		go Generator(limit, channel, quit, wg1, i)
	}

	for i := 0; i < int(count); i++ {
		quit <- true
		arr = append(arr, <-channel)
		// fmt.Printf("[%d] %d\n", i+1, x)
	}
	close(quit)
	wg1.Wait()
	close(channel)
	return arr
}

// Generator создаёт случайные числа от 0 до лимита и передаёт в канал
func Generator(limit int, channel chan int, quit chan bool, wg1 *sync.WaitGroup, index int) {
	defer wg1.Done()
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
