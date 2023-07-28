package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

type File struct {
	Type string
	Name string
	Size int64
}

type ServerOptions struct {
	Ip   string
	Port string
}

// type FileError struct {
// 	typeError int
// 	textError string
// }

// HomeHandler отправляет заголовок и вызывает открытие главной страницы
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var v = struct {
		Host string
	}{
		r.Host,
	}
	templ, _ := template.ParseFiles("./static/index.html")
	templ.Execute(w, &v)
}

//DirHandler получает адрес и тип сортировки, сортирует и отправляет массив объектов в формате json
func DirHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	root, sortType := r.URL.Query().Get("root"), r.URL.Query().Get("sortType")
	fmt.Printf("Получен %s\n", root)

	files, err := listDirByReadDir(root)
	if err != nil {
		fmt.Printf("Ошибка: \n", err)
		w.Header().Add("error", "fileNotExist")
		return
	}

	files = SortFiles(files, sortType)
	json_data, err := json.Marshal(files)
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(json_data)
	if err != nil {
		fmt.Println(err)
	}
}

//getIpPort читает и возвращает ip и port из файла config.json
func getIpPort() (string, string, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return "", "", fmt.Errorf("Ошибка открытия конфига сервера: ", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("Ошбика чтения данных из файла: ", err)
	}

	var serverOptions ServerOptions
	err = json.Unmarshal(data, &serverOptions)
	if err != nil {
		return "", "", fmt.Errorf("Ошибка чтения данных из файла: ", err)
	}
	return serverOptions.Ip, serverOptions.Port, nil
}

// listDirByReadDir возвращает типы, названия подпапок, файлов и их размеры по переданному адресу
func listDirByReadDir(path string) ([]File, error) {
	var allFiles []File
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return allFiles, fmt.Errorf("ошибка с чтением директории: %s", err)
	}
	for _, file := range lst {
		if file.IsDir() { // Обработка подпапки и её внутренностей
			folderSize, err := GetFolderSize(path + "/" + file.Name())
			fmt.Printf("%s : %d\n", file.Name(), folderSize)
			if err != nil {
				fmt.Printf("ошибка с чтением директории: %s", err)
				continue
			}
			allFiles = append(allFiles, File{"folder", fmt.Sprintf("%s/%s", path, file.Name()), folderSize})
		} else { // Обработка файла
			sizeOfFile, err := GetFileSize(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				fmt.Printf("ошибка при чтении размера файла: %s", err)
				continue
			}
			allFiles = append(allFiles, File{"file", fmt.Sprintf("%s/%s", path, file.Name()), sizeOfFile})
		}
	}
	return allFiles, nil
}

// GetFolderSize возвращает размер папки, по указанному адресу
func GetFolderSize(path string) (int64, error) {
	var folderSize int64 = 0
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return folderSize, fmt.Errorf("ошибка с чтением директории: %s", err)
	}
	for _, file := range lst {
		if file.IsDir() { // Обработка подпапки и её внутренностей
			subFolderSize, err := GetFolderSize(path + "/" + file.Name())
			if err != nil {
				fmt.Printf("ошибка при чтении размера файла: %s", err)
				continue
			}
			folderSize += subFolderSize
		} else { // Обработка файла
			sizeOfFile, err := GetFileSize(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				fmt.Printf("ошибка при чтении размера файла: %s", err)
				continue
			}
			folderSize += sizeOfFile
		}
	}
	return folderSize, nil
}

// GetFileSize возвращает размер файла, по указанному адресу
func GetFileSize(root string) (int64, error) {
	fInfo, err := os.Stat(root)
	if err != nil {
		return -1, fmt.Errorf("ошибка с поиском директории: %s", err)
	}
	return fInfo.Size(), nil
}

//SortFiles сортирует файлы по возрастанию по умолчанию или убываню и возвращает их
func SortFiles(files []File, sortType string) []File {
	if strings.EqualFold(sortType, "asc") || sortType == "" {
		sort.Slice(files, func(i, j int) (less bool) {
			return files[i].Size < files[j].Size
		})
	} else {
		sort.Slice(files, func(i, j int) (less bool) {
			return files[i].Size > files[j].Size
		})
	}
	return files
}

func main() {
	ip, port, err := getIpPort()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Handle("/static/styles/", http.StripPrefix("/static/styles/", http.FileServer(http.Dir("./static/styles/"))))
	http.Handle("/static/script/", http.StripPrefix("/static/script/", http.FileServer(http.Dir("./static/script/"))))
	http.Handle("/static/images/", http.StripPrefix("/static/images/", http.FileServer(http.Dir("./static/images/"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/dir", DirHandler)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil); err != nil {
		fmt.Println(err)
	}
}
