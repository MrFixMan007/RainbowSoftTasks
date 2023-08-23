package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type responseFiles struct {
	Time  string
	Files []File
}

type File struct {
	Type string
	Name string
	Size int64
}

type ServerOptions struct {
	Port string
}

// HomeHandler отправляет заголовок и вызывает открытие главной страницы
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var v = struct {
		Host string
	}{
		r.Host,
	}
	templ, _ := template.ParseFiles("./public/index.html")
	templ.Execute(w, &v)
}

//DirHandler получает адрес и тип сортировки, сортирует и отправляет массив объектов в формате json
func DirHandler(w http.ResponseWriter, r *http.Request) {
	//Узнаём время
	start := time.Now()

	//Ставим заголовок тип контента на json
	w.Header().Set("Content-Type", "application/json")

	//получаем входные параметры из запроса
	root, sortType := r.URL.Query().Get("root"), r.URL.Query().Get("sortType")
	fmt.Printf("Получен %s\n", root)

	//Подсчитываем файлы/папки по адресу
	files, err := listDirByReadDir(root)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("error", "fileNotExist")
		return
	}

	//Сортировка по размеру
	files = SortFiles(files, sortType)

	//Вычисляем сколько времени заняло
	elapsedTime := float64(time.Since(start) / time.Millisecond)
	elapsedTimeStr := fmt.Sprintf("%0.1f секунд(ы)", float64(elapsedTime/1000))

	//Маршалим ответ в вид json
	response := responseFiles{Time: elapsedTimeStr, Files: files}
	json_data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	//Отправляем ответ
	_, err = w.Write(json_data)
	if err != nil {
		fmt.Println(err)
	}

	//получаем ip
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	go addStatistic(root, elapsedTime, ip)
}

//addStatistic отправляет Put-запрос на сервер Apache, который записывает статистику в БД
func addStatistic(root string, elapsedTime float64, ip string) error {

	//адрес на php-скрипт сервера Apache
	url := fmt.Sprintf("http://%s:80/stat.php", ip)

	//создаём переменную, которую передадим
	values := map[string]string{"root": root, "elapsed_time": fmt.Sprintf("%0.0f", elapsedTime)}

	//маршалим в формат json
	json_data, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err)
	}

	//создаём Put-запрос, в который передаём адрес и json-данные
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
	}

	//устанавливает заголовок
	req.Header.Set("Content-Type", "application/json")

	//создаём клиент
	client := &http.Client{}

	//отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(resp, err)
	}
	defer resp.Body.Close()

	//печатаем ответ
	fmt.Println(resp)
	return nil
}

//getIpPort читает и возвращает ip и port из файла config.json
func getPort() (string, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return "", fmt.Errorf("Ошибка открытия конфига сервера: ", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Ошибка чтения данных из файла: ", err)
	}

	var serverPort ServerOptions
	err = json.Unmarshal(data, &serverPort)
	if err != nil {
		return "", fmt.Errorf("Ошибка чтения данных из файла: ", err)
	}
	return serverPort.Port, nil
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
			subFolderSize, err := GetFolderSize(path + "/" + file.Name()) //получаем размер подпапки
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

//main получает необходимые конфигурационные данные, ставит обработчики
func main() {
	port, err := getPort()
	if err != nil {
		fmt.Println(err)
		return
	}
	// http.Handle("/public/styles/", http.StripPrefix("/public/styles/", http.FileServer(http.Dir("./public/styles/"))))
	http.Handle("/public/dist/", http.StripPrefix("/public/dist/", http.FileServer(http.Dir("./public/dist/"))))
	// http.Handle("/public/images/", http.StripPrefix("/public/images/", http.FileServer(http.Dir("./public/images/"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/dir", DirHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Println(err)
	}
}
