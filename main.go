package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./src/templates/base.html", "./src/templates/home.html"))
	t.ExecuteTemplate(w, "base", nil)
}

func main() {
	cssFolder := http.FileServer(http.Dir("src/css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFolder))

	imgFolder := http.FileServer(http.Dir("src/img"))
	http.Handle("/img/", http.StripPrefix("/img/", imgFolder))

	jsFolder := http.FileServer(http.Dir("src/js"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFolder))

	scssFolder := http.FileServer(http.Dir("src/scss"))
	http.Handle("/scss/", http.StripPrefix("/scss/", scssFolder))

	vendorFolder := http.FileServer(http.Dir("src/vendor"))
	http.Handle("/vendor/", http.StripPrefix("/vendor/", vendorFolder))

	fmt.Printf("Escuchando...")
	http.HandleFunc("/", handleHomePage)
	http.ListenAndServe(":8080", nil)
}
