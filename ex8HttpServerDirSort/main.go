package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

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

}

func getIpPort() (string, string, error) {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Errorf("Ошибка открытия конфига сервера: ", err)
		return "", "", err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Errorf("Ошбика чтения данных из файла: ", err)
		return "", "", err
	}

	var serverOptions ServerOptions
	err = json.Unmarshal(data, &serverOptions)
	if err != nil {
		fmt.Errorf("Ошибка размаршалинга файла: ", err)
		return "", "", err
	}
	return serverOptions.Ip, serverOptions.Port, nil
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
