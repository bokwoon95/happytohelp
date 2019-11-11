package main

import (
	"html/template"
	"log"
	"net/http"
)

func landing(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "landing.html")
}

func category(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "category.html")
}

func disclosure(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("disclosure.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	case "POST":
		w.Write([]byte("hey ho"))
	}
}

func main() {
	http.HandleFunc("/", landing)
	http.HandleFunc("/category", category)
	http.HandleFunc("/disclosure", disclosure)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
