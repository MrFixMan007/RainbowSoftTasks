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

// Printer создаёт горутины, передаёт им параметры. Выведет столько чисел, сколько было передано.
func Printer(flow int64, limit int64, count int64, wg *sync.WaitGroup) error {
	defer wg.Done()
	channel := make(chan int)
	wg1 := &sync.WaitGroup{}
	for i := 0; (int64)(i) < flow; i++ {
		wg1.Add(1)
		go Generator(limit, channel, wg1, i)
		//fmt.Printf("Создал %d-ю горутину\n", i+1)
	}

	for i := 0; i < int(count); i++ {
		x := <-channel
		fmt.Printf("[%d] %d\n", i+1, x)
	}
	wg1.Add(-int(flow))
	return nil
}

// Generator создаёт случайные числа от 0 до лимита и передаёт в канал
func Generator(limit int64, channel chan int, wg1 *sync.WaitGroup, index int) {
	//TODO Зацикливание
	for {
		source := rand.NewSource(time.Now().UnixNano())
		r2 := rand.New(source)
		channel <- r2.Intn(int(limit))
	}
}
