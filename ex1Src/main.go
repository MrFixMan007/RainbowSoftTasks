package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	//Флаги
	src, dst, err := flagParsing()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	//Читаем файл
	srcs, err := ReadFile(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//Выполняем запросы
	for i, v := range srcs {
		body, err := RequestGets(v)
		if err != nil {
			fmt.Println(err)
		} else {
			err = WriteFile(body, dst, i)
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}
		}
	}
	fmt.Printf("Время выполнения: %0.3f секунд\n", float64(time.Since(start)/time.Millisecond)/1000)
}

// flagParsing проверяем и возвращаем переданные параметры
func flagParsing() (string, string, error) {
	src := flag.String("src", "", "Адрес на файл с источниками должен быть указан в формате: /home/...")
	dst := flag.String("dst", "", "Адрес на файл с местом сохранения должен быть указан в формате: /home/...")
	flag.Parse()

	if *src == "" || *dst == "" {
		return "", "", errors.New("")
	}
	return *src, *dst, nil
}

// ReadFile получаем текст файла по адресу
func ReadFile(src string) ([]string, error) {
	var strs []string
	file, err := os.Open(src)
	if err != nil {
		return strs, fmt.Errorf("ошибка при открытии файла: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return strs, fmt.Errorf("ошибка при чтении файла: %s", err)
	}
	return strs, nil
}

// RequestGets получаем ответы на запросы по переданной ссылке
func RequestGets(v string) ([]byte, error) {
	resp, err := http.Get(v)
	var body []byte
	if err != nil {
		return body, fmt.Errorf("ошибка с запросом: %s", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, fmt.Errorf("ошибка с чтением ответа: %s", err)
	}
	return body, nil
}

// WriteFile записываем по адресу текст в файл
func WriteFile(body []byte, dst string, i int) error {
	file, err := os.Create(fmt.Sprintf("%s%d", dst, i))
	if err != nil {
		return fmt.Errorf("ошибка с созданием файла: %s", err)
	}
	defer file.Close()

	n, err := file.WriteString(string(body))
	if err != nil {
		return fmt.Errorf("ошибка с записью в файл: %s (%d)", err, n)
	}
	return nil
}
