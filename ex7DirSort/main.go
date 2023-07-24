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

func main() {
	type File struct {
		Type string
		Name string
		Size int64
	}
	//Флаги
	root, sortType, err := flagParsing()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	types, names, values, err := listDirByReadDir(root)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var files = []File{}
	for i, v := range values {
		files = append(files, File{types[i], names[i], v})
	}
	if strings.EqualFold(sortType, "asc") {
		fmt.Println("Сортировка в порядке возрастания")
		sort.Slice(files, func(i, j int) (less bool) {
			return files[i].Size < files[j].Size
		})
	} else {
		fmt.Println("Сортировка в порядке убывания")
		sort.Slice(files, func(i, j int) (less bool) {
			return files[i].Size > files[j].Size
		})
	}

	for i, _ := range files {
		fmt.Printf("[%d] %s %s %d\n", i, files[i].Type, files[i].Name, files[i].Size)
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

// listDirByReadDir возвращает названия подпапок, файлов и из размеры по переданному адресу
func listDirByReadDir(path string) ([]string, []string, []int64, error) {
	var types []string
	var names []string
	var sizes []int64
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return types, names, sizes, fmt.Errorf("ошибка с чтением директории: %s", err)
	}
	for _, val := range lst {
		if val.IsDir() { // Обработка подпапки и её внутренностей
			types1, names1, sizes1, err := listDirByReadDir(fmt.Sprintf("%s/%s", path, val.Name())) //Рекурсивный вызов функции
			if err != nil {
				fmt.Printf("ошибка с чтением директории: %s", err)
				continue
			}
			value := 0
			for _, v := range sizes1 { //Подсчитываем размер папки
				value += int(v)
			}

			types = append(types, "dir")
			names = append(names, fmt.Sprintf("%s/%s", path, val.Name())) //Добавляем подпапку
			sizes = append(sizes, int64(value))                           //Добавляем размер подпапки

			types = append(types, types1...)
			names = append(names, names1...) //Добавляем подфайлы
			sizes = append(sizes, sizes1...) //Добавляем размеры подфайлов
			//fmt.Printf("[%s]\n", fmt.Sprintf("%s/%s", path, val.Name()))
		} else { // Обработка файла
			value, err := GetSize(fmt.Sprintf("%s/%s", path, val.Name()))
			if err != nil {
				fmt.Printf("ошибка при чтении размера файла: %s", err)
				continue
			}
			types = append(types, "file")
			sizes = append(sizes, value)                                  //Добавляем файлы
			names = append(names, fmt.Sprintf("%s/%s", path, val.Name())) //Добавляем размеры файлов
		}
	}
	return types, names, sizes, nil
}
