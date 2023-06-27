package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}

	flow, limit, count, err := flagParsing()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	wg.Add(1)
	go Printer(flow, limit, count, wg)
	wg.Wait()
}

// flagParsing проверяем и возвращаем переданные параметры
func flagParsing() (int64, int64, int64, error) {
	flow := flag.Int64("flow", 0, "Количество потоков должно быть указано в формате натурального числа большего 0")
	limit := flag.Int64("limit", 0, "Верхняя граница случайных чисел должна быть указана в формате натурального числа")
	count := flag.Int64("count", 0, "Количество чисел должно быть указано в формате натурального числа большего 0")
	flag.Parse()

	if *flow <= 0 || *limit <= 0 || *count <= 0 {
		return 0, 0, 0, errors.New("ошибка с обработкой флагов")
	}
	return *flow, *limit, *count, nil
}

// Printer создаёт горутины, передаёт им параметры. Выведет столько чисел, количество которых было передано.
func Printer(flow int64, limit int64, count int64, wg *sync.WaitGroup) {
	defer wg.Done()
	channel := make(chan int)
	quit := make(chan bool)
	wg1 := &sync.WaitGroup{}
	for i := 0; (int64)(i) < flow; i++ {
		wg1.Add(1)
		go Generator(limit, channel, quit, wg1, i)
	}

	for i := 0; i < int(count); i++ {
		quit <- true
		x := <-channel
		fmt.Printf("[%d] %d\n", i+1, x)
	}
	close(quit)
	wg1.Wait()
	close(channel)
}

// Generator создаёт случайные числа от 0 до лимита и передаёт в канал
func Generator(limit int64, channel chan int, quit chan bool, wg1 *sync.WaitGroup, index int) {
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
