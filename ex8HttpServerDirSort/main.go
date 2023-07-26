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

// HomeHandler отправляет заголовок и вызывает открытие главной страницы
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var v = struct {
		Host string
	}{
		r.Host,
	}
	templ, _ := template.ParseFiles("index.html")
	templ.Execute(w, &v)
}

func DirHandler(w http.ResponseWriter, r *http.Request) {

	root, sortType := r.URL.Query().Get("root"), r.URL.Query().Get("sortType")

	files, err := listDirByReadDir(root)
	if err != nil {
		_, err := w.Write([]byte(fmt.Sprintf("%s", err)))
		if err != nil {
			fmt.Println(err)
		}
	}

	files = SortFiles(files, sortType) //сортировка

	json_data, err := json.Marshal(files)
	if err != nil {
		fmt.Println(err)
	}
	var file1 []File
	err = json.Unmarshal(json_data, &file1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(file1[1])
	}

	num, err := w.Write(json_data)
	if err != nil {
		fmt.Println(num, err)
	}
	// var responseInBytes []byte
	// for i := range files {
	// 	responseInBytes = append(responseInBytes, []byte(fmt.Sprintf("%s %s %d байт(а)\n", files[i].Type, files[i].Name, files[i].Size))...)
	// 	fmt.Printf("[%d] %s %s %d байт(а)\n", i+1, files[i].Type, files[i].Name, files[i].Size)
	// }

	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(files); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("JSON Закодирован")
	// }
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
		return "", "", fmt.Errorf("Ошбика чтения данных из файла: ", err)
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

//addToListDirWithSize считает размер папки и добавляет в массив данные папки через указатель
// func addToListDirWithSize(allFiles *[]File, filesInSubDir []File, path string, name string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	sizeOfDirectory := 0
// 	for i := range filesInSubDir {
// 		sizeOfDirectory += int(filesInSubDir[i].Size)
// 	}

// 	*allFiles = append(*allFiles, File{"dir", fmt.Sprintf("%s/%s", path, name), int64(sizeOfDirectory)})

// 	*allFiles = append(*allFiles, filesInSubDir...)
// }

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

func main() {
	ip, port, err := getIpPort()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles/"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/dir", DirHandler)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil); err != nil {
		fmt.Println(err)
	}
}
