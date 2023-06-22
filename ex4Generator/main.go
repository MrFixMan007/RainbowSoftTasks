package main

import (
	"errors"
	"flag"
	"os"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	channel := make(chan int)

	flow, limit, count, err := flagParsing()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	go Printer(count, channel, wg)

	for i := 0; (int64)(i) < flow; i++ {
		wg.Add(1)
		go Generator(limit, wg)
	}
}

// flagParsing проверяем и возвращаем переданные параметры
func flagParsing() (int64, int64, int64, error) {
	flow := flag.Int64("flow", 0, "Количество потоков должно быть указано в формате натурального числа большего 0")
	limit := flag.Int64("limit", 0, "Верхняя граница случайных чисел должна быть указана в формате натурального числа")
	count := flag.Int64("count", 0, "Количество чисел должно быть указано в формате натурального числа большего 0")
	flag.Parse()

	if *flow <= 0 || *limit < 0 || *count <= 0 {
		return 0, 0, 0, errors.New("")
	}
	return *flow, *limit, *count, nil
}

// WriteFile записываем по адресу текст в файл
func Printer(count int64, c chan int, wg *sync.WaitGroup) error {
	// i := 0;
	return nil
}

//
func Generator(limit int64, wg *sync.WaitGroup) error {
	return nil
}
