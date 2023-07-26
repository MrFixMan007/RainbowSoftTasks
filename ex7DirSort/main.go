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
	"sync"
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

	files = SortFiles(files, sortType)

	for i := range files {
		fmt.Printf("[%d] %s %s %d байт(а)\n", i+1, files[i].Type, files[i].Name, files[i].Size)
	}
}

// flagParsing проверяем и возвращаем переданные параметры
func flagParsing() (string, string, error) {
	root := flag.String("root", "", "Адрес на директорию должен быть указан в формате: /home/...")
	sortType := flag.String("sortType", "", "Укажите тип сортировки: ASC - по возрастаню, DESC - по убыванию. По умолчанию ASC")
	flag.Parse()

	matched, _ := regexp.MatchString("^/home/", *root) //Проверяем начало адреса
	if *root == "" || (*sortType != "" && !strings.EqualFold(*sortType, "asc") && !strings.EqualFold(*sortType, "desc")) || !matched {
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
	wg := &sync.WaitGroup{}
	var allFiles []File
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return allFiles, fmt.Errorf("ошибка с чтением директории: %s", err)
	}
	for _, file := range lst {
		if file.IsDir() { // Обработка подпапки и её внутренностей
			filesInSubDir, err := listDirByReadDir(fmt.Sprintf("%s/%s", path, file.Name())) //Рекурсивный обход папки
			if err != nil {
				fmt.Printf("ошибка с чтением директории: %s", err)
				continue
			}
			wg.Add(1)
			go addToListDirWithSize(&allFiles, filesInSubDir, path, file.Name(), wg)
		} else { // Обработка файла
			sizeOfFile, err := GetSize(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				fmt.Printf("ошибка при чтении размера файла: %s", err)
				continue
			}
			allFiles = append(allFiles, File{"file", fmt.Sprintf("%s/%s", path, file.Name()), sizeOfFile})
		}
	}
	wg.Wait()
	return allFiles, nil
}

//addToListDirWithSize считает размер папки и добавляет в массив данные папки через указатель
func addToListDirWithSize(allFiles *[]File, filesInSubDir []File, path string, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	sizeOfDirectory := 0
	for i := range filesInSubDir {
		sizeOfDirectory += int(filesInSubDir[i].Size)
	}

	*allFiles = append(*allFiles, File{"dir", fmt.Sprintf("%s/%s", path, name), int64(sizeOfDirectory)})

	*allFiles = append(*allFiles, filesInSubDir...)
}

//SortFiles сортирует файлы по возрастанию по умолчанию или убываню и возвращает их
func SortFiles(files []File, sortType string) []File {
	if strings.EqualFold(sortType, "asc") || sortType == "" {
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
	return files
}
