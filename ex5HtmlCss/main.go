package main

import (
	"html/template"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("index.html")
	templ.Execute(w, nil)
	// fmt.Fprintf(w, StyleAdress)
}

func start() {
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles/"))))
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":4434", nil)
}

func main() {
	start()
}
