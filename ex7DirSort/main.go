package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

type File struct {
	Type string
	Name string
	Size int64
}

func main() {

	//Флаги
	root, sortType, err := flagParsing()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	files, err := listDirByReadDir(root)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//сортировка
	if strings.EqualFold(sortType, "asc") {
		fmt.Println("Сортировка в порядке возрастания")
		sort.Slice(files, func(i, j int) (less bool) {
			return files[i].Size < files[j].Size
		})
	} else {
		fmt.Println("Сортировка в порядке убывания:")
		sort.Slice(files, func(i, j int) (less bool) {
			return files[i].Size > files[j].Size
		})
	}

	for i := range files {
		fmt.Printf("[%d] %s %s %d байт(а)\n", i+1, files[i].Type, files[i].Name, files[i].Size)
	}
}

// flagParsing проверяем и возвращаем переданные параметры
func flagParsing() (string, string, error) {
	root := flag.String("root", "", "Адрес на директорию должен быть указан в формате: /home/...")
	sortType := flag.String("sortType", "", "Укажите тип сортировки: ASC - по возрастаню, DESC - по убыванию")
	flag.Parse()

	matched, _ := regexp.MatchString("^/home/", *root) //Проверяем начало адреса
	if *root == "" || (!strings.EqualFold(*sortType, "asc") && !strings.EqualFold(*sortType, "desc")) || !matched {
		return "", "", errors.New("ошибка с обработкой флагов")
	}
	return *root, *sortType, nil
}

// GetSize возвращает размер файла, по указанному адресу
func GetSize(root string) (int64, error) {
	fInfo, err := os.Stat(root)
	if err != nil {
		return -1, fmt.Errorf("ошибка с поиском директории: %s", err)
	}
	return fInfo.Size(), nil
}

// listDirByReadDir возвращает типы, названия подпапок, файлов и их размеры по переданному адресу
func listDirByReadDir(path string) ([]File, error) {
	var files []File
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return files, fmt.Errorf("ошибка с чтением директории: %s", err)
	}
	for _, val := range lst {
		if val.IsDir() { // Обработка подпапки и её внутренностей
			files1, err := listDirByReadDir(fmt.Sprintf("%s/%s", path, val.Name())) //Рекурсивный вызов функции
			if err != nil {
				fmt.Printf("ошибка с чтением директории: %s", err)
				continue
			}
			value := 0
			for i := range files1 { //Подсчитываем размер папки
				value += int(files1[i].Size)
			}

			files = append(files, File{"dir", fmt.Sprintf("%s/%s", path, val.Name()), int64(value)})

			files = append(files, files1...)
		} else { // Обработка файла
			value, err := GetSize(fmt.Sprintf("%s/%s", path, val.Name()))
			if err != nil {
				fmt.Printf("ошибка при чтении размера файла: %s", err)
				continue
			}
			files = append(files, File{"file", fmt.Sprintf("%s/%s", path, val.Name()), value})
		}
	}
	return files, nil
}
