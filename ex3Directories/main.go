package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	var adresses []string
	//Флаги
	root, sizeLimit, err := flagParsing()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	names, values, err := listDirByReadDir(root)
	if err != nil {
		fmt.Println("Не удалось") //TODO
		os.Exit(1)                //TODO
	}

	for i, v := range values {
		if v > sizeLimit {
			fmt.Printf("Размер файла %s = %d\n", names[i], v)
			adresses = append(adresses, fmt.Sprintf("%s/%s", root, names[i]))
		}
	}
	err = WriteFile(adresses, root)
	if err != nil {
		fmt.Println("Не удалось") //TODO
		os.Exit(1)                //TODO
	}
	fmt.Printf("Время выполнения: %0.3f секунд\n", float64(time.Since(start)/time.Millisecond)/1000)
}

// flagParsing проверяем и возвращаем переданные параметры
func flagParsing() (string, int64, error) {
	root := flag.String("root", "", "Адрес на директорию должен быть указан в формате: /home/...")
	sizeLimit := flag.Int64("sizeLimit", 0, "Максимальный размер (в Кб) должен быть указан в формате числа")
	flag.Parse()

	if *root == "" || *sizeLimit == 0 {
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
		return fmt.Errorf("ошибка с записью в файл: %s (%d)", err, n)
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

func listDirByReadDir(path string) ([]string, []int64, error) {
	var body []int64
	var names []string
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return names, body, fmt.Errorf("ошибка с поиском директории: %s", err)
	}
	for _, val := range lst {
		if val.IsDir() {
			// names1, values1, err1 := listDirByReadDir(fmt.Sprintf("%s/%s", path, val.Name()))
			// if err1 != nil {
			// 	fmt.Println("Не удалось") //TODO
			// 	os.Exit(1)
			// }
			// body = append(body, value)
			// names = append(names, val.Name())
			//TODO
			//с файлами работает
			//TODO
		} else {
			value, err := GetSize(fmt.Sprintf("%s/%s", path, val.Name()))
			if err != nil {
				fmt.Printf("ошибка с при чтении размера файла: %s", err)
				continue
			}
			body = append(body, value)
			names = append(names, val.Name())
		}
	}
	return names, body, nil
}
