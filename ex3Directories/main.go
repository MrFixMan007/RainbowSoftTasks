package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	var adresses []string
	//Флаги
	root, sizeLimit, err := flagParsing()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	names, values, err := listDirByReadDir(root)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for i, v := range values {
		if v > sizeLimit {
			fmt.Printf("[%d] %s с размером = %d > %d\n", i, names[i], v, sizeLimit)
			adresses = append(adresses, names[i])
		} else {
			fmt.Printf("[%d] %s с размером = %d\n", i, names[i], v)
		}
	}
	err = WriteFile(adresses, root)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

// flagParsing проверяем и возвращаем переданные параметры
func flagParsing() (string, int64, error) {
	root := flag.String("root", "", "Адрес на директорию должен быть указан в формате: /home/...")
	sizeLimit := flag.Int64("sizeLimit", 0, "Максимальный размер (в байтах) должен быть указан в формате натурального числа")
	flag.Parse()

	matched, _ := regexp.MatchString("^/home/", *root) //Проверяем начало
	if *root == "" || *sizeLimit < 0 || !matched {
		return "", 0, errors.New("")
	}
	return *root, *sizeLimit, nil
}

// WriteFile записываем по адресу текст в файл
func WriteFile(body []string, src string) error {
	file, err := os.Create(fmt.Sprintf("%s/out", src))
	if err != nil {
		return fmt.Errorf("ошибка с созданием файла: %s", err)
	}
	defer file.Close()

	n, err := file.WriteString(strings.Join(body, "\n"))
	if err != nil {
		fmt.Printf("ошибка с записью в файл: %s (%d)", err, n)
	}
	return nil
}

// GetSize возвращает размер файла, по указанному адресу
func GetSize(root string) (int64, error) {
	fInfo, err := os.Stat(root)
	if err != nil {
		return -1, fmt.Errorf("ошибка с поиском директории: %s", err)
	}
	return fInfo.Size(), nil
}

// listDirByReadDir возвращает названия подпапок, файлов и из размеры по переданному адресу
func listDirByReadDir(path string) ([]string, []int64, error) {
	var sizes []int64
	var names []string
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return names, sizes, fmt.Errorf("ошибка с чтением директории: %s", err)
	}
	for _, val := range lst {
		if val.IsDir() { // Обработка подпапки и её внутренностей
			names1, sizes1, err := listDirByReadDir(fmt.Sprintf("%s/%s", path, val.Name())) //Рекурсивный вызов функции
			if err != nil {
				fmt.Printf("ошибка с чтением директории: %s", err)
				continue
			}
			value := 0
			for _, v := range sizes1 { //Подсчитываем размер папки
				value += int(v)
			}

			names = append(names, fmt.Sprintf("%s/%s", path, val.Name())) //Добавляем подпапку
			sizes = append(sizes, int64(value))                           //Добавляем размер подпапки

			names = append(names, names1...) //Добавляем подфайлы
			sizes = append(sizes, sizes1...) //Добавляем размеры подфайлов
			fmt.Printf("[%s]\n", fmt.Sprintf("%s/%s", path, val.Name()))
		} else { // Обработка файла
			value, err := GetSize(fmt.Sprintf("%s/%s", path, val.Name()))
			if err != nil {
				fmt.Printf("ошибка при чтении размера файла: %s", err)
				continue
			}
			sizes = append(sizes, value)                                  //Добавляем файлы
			names = append(names, fmt.Sprintf("%s/%s", path, val.Name())) //Добавляем размеры файлов
		}
	}
	return names, sizes, nil
}
